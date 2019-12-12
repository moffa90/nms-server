package serialWorker

import (
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils/usb"
	"github.com/moffa90/nms-server/worker"
	log "github.com/Sirupsen/logrus"
	"sync"
)

var queue *worker.WorkQueue
var once sync.Once
var lastResponse map[string]string

var QuitChan = make(chan bool)

func GetQueue() *worker.WorkQueue {
	once.Do(func() {
		queue = &worker.WorkQueue{
			Pending: make(chan worker.WorkRequest, 100),
		}
	})

	return queue
}

func GetLastResponse() map[string]string{
	return lastResponse
}

func StartWorker(){
	//insertUpdateWorkPerHardware()
	GetQueue()
	//log.Printf("%#V", queue.Pending)
	for{
		//log.Printf("Left %d works",len(queue.Pending))
		work, ok := <-queue.Pending
		if ok {
			// Receive a work request.
			switch work.GetAction() {
			case constants.SerialWorkRequestUpdateDevice:
				updateFullDevice(work.(*worker.SerialWorkRequest).Device)
				if work.GetWg() != nil {
					work.GetWg().Done()
				}
				break
			case constants.SerialWorkRequestUpdateSensorsDevice:
				updateSensorsDevice(work.(*worker.SerialWorkRequest).Device)
				if work.GetWg() != nil {
					work.GetWg().Done()
				}
				break
			case constants.SerialWorkRequestUpdateInfoDevice:
				updateInfoDevice(work.(*worker.SerialWorkRequest).Device)
				if work.GetWg() != nil {
					work.GetWg().Done()
				}
				break
			case constants.SerialWorkRequestExecCommand:
				lastResponse = execCommand(work.(*worker.SerialWorkRequest).Device, work.(*worker.SerialWorkRequest).RawCommand)
				if work.GetWg() != nil {
					work.GetWg().Done()
				}
				break
			}

		} else {
			log.Printf("%#v", work)
		}
	}
}

func updateFullDevice(h models.Hardware){
	//start := time.Now()
	if data := usb.GetFullData(h.Serial); data != nil {
		if err := h.Update(data["info"], db.Shared); err != nil {
			log.Error(err)
		}

		for cat,fields := range data {
			if cat != "info"{
				for key,value := range fields {
					if e := models.NewHardwareCharacteristics(h.Id, h.Serial, key, value, cat, h).Save(db.Shared); e != nil {
						log.Error(e)
					}
				}
			}

		}
	}else{
		log.Errorf("%s not reached.", h.DevId)
	}

	//log.Printf("End update %s: %s", h.Serial, time.Since(start))
}

func updateSensorsDevice(h models.Hardware){
	//start := time.Now()
	if data := usb.GetSensorsData(h.Serial); data != nil {
		for cat,fields := range data {
			for key,value := range fields {
				if e := models.NewHardwareCharacteristics(h.Id, h.Serial, key, value, cat, h).Save(db.Shared); e != nil {
					log.Error(e)
				}
			}
		}
	}else{
		log.Errorf("%s not reached.", h.DevId)
	}

	//log.Printf("End update %s: %s", h.Serial, time.Since(start))
}

func updateInfoDevice(h models.Hardware){
	//start := time.Now()
	if data := usb.GetInfoData(h.Serial); data != nil {
		if err := h.Update(data["info"], db.Shared); err != nil {
			log.Error(err)
		}
	}else{
		log.Errorf("%s not reached.", h.DevId)
	}

	//log.Printf("End update %s: %s", h.Serial, time.Since(start))
}

func execCommand(h models.Hardware, c string) map[string]string{
	//start := time.Now()
	rawResponse, response := usb.ExecCommand(h.Serial, c)
	log.Printf(rawResponse)
	return response
	//log.Printf("End execution of %s for %s: %s", c, h.Serial, time.Since(start))
}

func insertUpdateWorkPerHardware(){

	hardware, err := models.GetHardware(db.Shared)
	if err != nil{
		log.Error(err)
	}

	for _,h := range hardware {
		work := worker.NewSerialWorkRequest(h, constants.SerialWorkRequestUpdateDevice, "", nil)
		queue.Pending<- *work
	}
}