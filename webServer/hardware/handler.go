package hardware

import (
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/utils"
	"github.com/moffa90/triadNMS/utils/security"
	"github.com/moffa90/triadNMS/utils/usb"
	"github.com/moffa90/triadNMS/worker"
	"github.com/moffa90/triadNMS/worker/serialWorker"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	h, _ := models.GetHardware(db.Shared)
	var wg sync.WaitGroup
	wg.Add(len(h))
	for _, hard := range h{
		work := worker.NewSerialWorkRequest(hard, constants.SerialWorkRequestUpdateInfoDevice, "", &wg)
		serialWorker.GetQueue().AddWork(work)
	}
	wg.Wait()

	h, _ = models.GetHardware(db.Shared)
	infoStruct := struct{
		Active string
		Info map[string]string
		Errors map[string]string
		DetectedDevices []usb.USBDevice
		RegisteredDevices []models.Hardware
	}{
		"management",
		security.CookieGetInfo(w, req),
		nil,
		usb.GetDevices(),
		h,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_HARDWARE_PATH, infoStruct)
}

func InfoHandler(w http.ResponseWriter, req *http.Request) {

	var wg sync.WaitGroup
	wg.Add(1)
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	_ , alarmsParam:= req.URL.Query()["alarms"]
	_ , sensorsParam:= req.URL.Query()["sensors"]

	infoStruct := struct {
		Active 		string
		Info   		map[string]string
		Errors 		map[string]string
		HardInfo   	models.Hardware
		Sensors 	[]models.HardwareCharacteristics
		Controls	[]models.HardwareCharacteristics
		Alarms		[]models.HardwareCharacteristics
		Thresholds 	[]models.HardwareCharacteristics
		RFCalDet	[]models.HardwareCharacteristics
		ADC			[]models.HardwareCharacteristics
		Extra1		[]models.HardwareCharacteristics
		Extra2		[]models.HardwareCharacteristics
		ExtraName1	string
		ExtraName2	string
	} {
		"management",
		security.CookieGetInfo(w, req),
		nil,
		models.Hardware{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		"",
		"",
	}

	if serial != "" {
		if hardInfo, e := models.GetHardwareBySerial(db.Shared, serial); e != nil {
			log.Error(e)
			var work = worker.NewSerialWorkRequest(models.Hardware{Serial: serial}, constants.SerialWorkRequestExecCommand, constants.INFO_COMMAND, &wg)
			serialWorker.GetQueue().AddWork(work)
			wg.Wait()
			data := serialWorker.GetLastResponse()
			data["Serial"] = serial
			if payload, err := json.Marshal(data); err == nil {
				log.Printf(string(payload))
				w.Header().Set("Content-Type", "application/json")
				if _, err = w.Write(payload); err != nil {
					log.Error(err)
				}
				return
			}
		}else{
			infoStruct.HardInfo = *hardInfo
		}

		if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			work := worker.NewSerialWorkRequest(infoStruct.HardInfo, constants.SerialWorkRequestUpdateSensorsDevice, "", &wg,)
			serialWorker.GetQueue().AddWork(work)
			wg.Wait()
			if alarmsParam {
				infoStruct.Alarms, _ = models.GetHardwareCharacteristicsByCat(serial,"alarms", db.Shared)
			}

			if sensorsParam {
				infoStruct.Sensors, _ = models.GetHardwareCharacteristicsByCat(serial,"sensors", db.Shared)
				infoStruct.ADC, _ = models.GetHardwareCharacteristicsByCat(serial,"ADC", db.Shared)
			}

			if payload, err := json.Marshal(infoStruct); err == nil {
				w.Header().Set("Content-Type", "application/json")
				if _, err = w.Write(payload); err != nil {
					log.Error(err)
				}
			}else {
				log.Error(err)
			}
			return
		}else{
			work := worker.NewSerialWorkRequest(infoStruct.HardInfo, constants.SerialWorkRequestUpdateDevice, "", &wg)
			serialWorker.GetQueue().AddWork(work)

			wg.Wait()
			infoStruct.Sensors, _ = models.GetHardwareCharacteristicsByCat(serial,"sensors", db.Shared)
			infoStruct.Alarms, _ = models.GetHardwareCharacteristicsByCat(serial,"alarms", db.Shared)
			infoStruct.Controls, _ = models.GetHardwareCharacteristicsByCat(serial,"controls", db.Shared)
			infoStruct.Thresholds, _ = models.GetHardwareCharacteristicsByCat(serial,"thresholds", db.Shared)
			infoStruct.RFCalDet, _ = models.GetHardwareCharacteristicsByCat(serial,"RFCalDet", db.Shared)
			infoStruct.ADC, _ = models.GetHardwareCharacteristicsByCat(serial,"ADC", db.Shared)
			infoStruct.Extra1, _ = models.GetHardwareCharacteristicsByCat(serial, "extra1", db.Shared)
			infoStruct.Extra2, _ = models.GetHardwareCharacteristicsByCat(serial, "extra2", db.Shared)

			switch infoStruct.HardInfo.ProductId {
			case constants.FOTX_200_REV1:
				infoStruct.ExtraName1 = "Laser APC"
				infoStruct.ExtraName2 = "Optical Power"
				break
			case constants.MOD_100_REV1:
				infoStruct.ExtraName1 = "RF Detector 2"
				infoStruct.ExtraName2 = "RF AGC"

				break
			case constants.FORX_200_REV1:
				infoStruct.ExtraName1 = "Optical Power"
				break

			case constants.MOD_200_REV1:

				break
			}
		}

		utils.RenderPage(w, constants.TEMPLATE_PAGE_HARDWARE_INFO_PATH, infoStruct)
		return
	}
	http.Redirect(w, req, "/management/hardware", http.StatusFound)
}

func SaveHandler(w http.ResponseWriter, req *http.Request) {

	serial := req.PostFormValue("serialWorker")
	productID := req.PostFormValue("product")
	deviceID := req.PostFormValue("device")
	address := req.PostFormValue("address")
	errors := make(map[string]string)

	if serial == "" {
		errors["0"] = "Invalid serialWorker reference"
	}

	if productID == "" {
		errors["1"] = "Product ID is empty"
	}

	if deviceID == "" {
		errors["2"] = "Device ID is empty"
	}

	if address == "" {
		errors["3"] = "Address slot is empty"
	}

	if len(errors) == 0 {
		addrInt, _ := strconv.ParseInt(strings.Replace(address, "x", "", -1), 16, 32)
		h := models.Hardware{
			DevId: deviceID,
			ProductId: productID,
			Address: int(addrInt),
			Serial: serial,
		}

		if id, err := models.CreateHardware(db.Shared, &h);  err == nil {
			http.Redirect(w, req, "/management/hardware", http.StatusFound)
			log.Printf("New Hardware created with ID: %s, %#v\n\r", id, h)
		}else {
			errors["4"] = err.Error()
		}

	}

	h, _ := models.GetHardware(db.Shared)
	infoStruct := struct{
		Active string
		Info map[string]string
		Errors map[string]string
		DetectedDevices []usb.USBDevice
		RegisteredDevices []models.Hardware
	}{
		"management",
		security.CookieGetInfo(w, req),
		errors,
		usb.GetDevices(),
		h,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_HARDWARE_PATH, infoStruct)
}

func DeleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serial := vars["serialWorker"]

	if err := db.Shared.Unscoped().Delete(models.Hardware{}, "id=?", serial).Error; err == nil{
		w.WriteHeader(200)
	}else{
		log.Println(err.Error())
	}

}

func ActionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	action := vars["action"]

	if h, e := models.GetHardwareBySerial(db.Shared, serial); e == nil {
		var wg sync.WaitGroup
		wg.Add(1)
		work := worker.NewSerialWorkRequest(*h, constants.SerialWorkRequestExecCommand, "", &wg)

		switch action {

		case "reset":
			work.RawCommand = constants.RESET_COMMAND + "\n"
			break
		case "save":
			work.RawCommand = constants.SAVE_COMMAND + "\n"
			break
		case "restore":
			work.RawCommand = constants.RESTORE_COMMAND + "\n"
			break
		case "clearalarms":
			work.RawCommand = constants.CLEAR_ALARMS_COMMAND + "\n"
			break

		}
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()
	}

	log.Println(serial + " " + action)

	if action == "reset" {
		time.Sleep(8 * time.Second)
	}

	http.Redirect(w, req, "/management/hardware/"+serial, http.StatusFound)
}

func EditValuesHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	control := req.PostFormValue("control")
	value := req.PostFormValue("value")

	if h, e := models.GetHardwareBySerial(db.Shared, serial); e == nil {
		var wg sync.WaitGroup
		wg.Add(1)
		work := worker.NewSerialWorkRequest(*h, constants.SerialWorkRequestExecCommand, "", &wg)
		switch strings.ToLower(control) {
		case "rfattenuator1":
			work.RawCommand = constants.SET_RF_ATTENUATOR_1_COMMAND + value + "\n"
			break
		case "rfattenuator2":
			work.RawCommand = constants.SET_RF_ATTENUATOR_2_COMMAND + value + "\n"
			break
		case "rfattenuator3":
			work.RawCommand = constants.SET_RF_ATTENUATOR_3_COMMAND + value + "\n"
			break
		case "agcmode":
			work.RawCommand = constants.SET_AGC_MODE_COMMAND + value + "\n"
			break
		case "rfsquelch":
			work.RawCommand = constants.SET_RF_SQUELCH_COMMAND + value + "\n"
			break
		case "rfswitch":
			work.RawCommand = constants.SET_RF_SWITCH_COMMAND + value + "\n"
			break

		default:
			break
		}
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()
		w.WriteHeader(200)
	}

}

func ThresholdsHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	productID := req.PostFormValue("productID")
	command := ""
	if h, e := models.GetHardwareBySerial(db.Shared, serial); e == nil {
		var wg sync.WaitGroup
		wg.Add(1)

		switch productID {
		case "CGW-FOTX-200-REV1":
			//writeAlarmThresholds: <TempHi> <TempLo> <RFDetHi> <RFDetLo> <OptPwrHi> <OptPwrLo> <LaserBiasHi> <LaserBiasLo>
			command = fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
				constants.WRITE_THRESHOLDS_COMMAND,
				req.PostFormValue("TemperatureAlarmHi"),
				req.PostFormValue("TemperatureAlarmLo"),
				req.PostFormValue("RFDetectorAlarmHi"),
				req.PostFormValue("RFDetectorAlarmLo"),
				req.PostFormValue("OpticalPowerAlarmHi"),
				req.PostFormValue("OpticalPowerAlarmLo"),
				req.PostFormValue("LaserBiasAlarmHi"),
				req.PostFormValue("LaserBiasAlarmLo"),
			)
			break
		case "CGW-MOD-100-REV1":
			//writeAlarmThresholds: <TempHi> <TempLo> <RFDet1Hi> <RFDet1Lo> <RFDet2Hi> <RFDet2Lo>
			command = fmt.Sprintf("%s %s %s %s %s %s %s",
				constants.WRITE_THRESHOLDS_COMMAND,
				req.PostFormValue("TemperatureAlarmHi"),
				req.PostFormValue("TemperatureAlarmLo"),
				req.PostFormValue("RFDetector1AlarmHi"),
				req.PostFormValue("RFDetector1AlarmLo"),
				req.PostFormValue("RFDetector2AlarmHi"),
				req.PostFormValue("RFDetector2AlarmLo"),
			)

			break
		case "CGW-FORX-200-REV1":
			//writeAlarmThresholds: <TempHi> <TempLo> <RFDet1Hi> <RFDet1Lo> <OptPwrHi> <OptPwrLo>
			command = fmt.Sprintf("%s %s %s %s %s %s %s",
				constants.WRITE_THRESHOLDS_COMMAND,
				req.PostFormValue("TemperatureAlarmHi"),
				req.PostFormValue("TemperatureAlarmLo"),
				req.PostFormValue("RFDetectorAlarmHi"),
				req.PostFormValue("RFDetectorAlarmLo"),
				req.PostFormValue("OpticalPowerAlarmHi"),
				req.PostFormValue("OpticalPowerAlarmLo"),
			)
			break

		case "CGW-MOD-200-REV1":
			//writeAlarmThresholds: <TempHi> <TempLo> <RFDetHi> <RFDetLo>
			command = fmt.Sprintf("%s %s %s %s %s",
				constants.WRITE_THRESHOLDS_COMMAND,
				req.PostFormValue("TemperatureAlarmHi"),
				req.PostFormValue("TemperatureAlarmLo"),
				req.PostFormValue("RFDetectorAlarmHi"),
				req.PostFormValue("RFDetectorAlarmLo"),
			)
			break
		}

		work := worker.NewSerialWorkRequest(*h, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(404)
}

func RfDetectorHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	command := ""

	if h, e := models.GetHardwareBySerial(db.Shared, serial); e == nil {
		var wg sync.WaitGroup
		wg.Add(1)
		//writeCalRFDet: <frequency> <intercept> <slope> <AlarmHi> <AlarmLo>
		command = fmt.Sprintf("%s %s %s %s %s %s",
			constants.WRITE_CAL_RF_DET_COMMAND,
			req.PostFormValue("Frequency"),
			req.PostFormValue("Intercept"),
			req.PostFormValue("Slope"),
			req.PostFormValue("AlarmThresholdHi"),
			req.PostFormValue("AlarmThresholdLo"),
		)
		work := worker.NewSerialWorkRequest(*h, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(404)
}

func ExtraHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	serial := vars["serialWorker"]
	productID := req.PostFormValue("productID")

	if h, e := models.GetHardwareBySerial(db.Shared, serial); e == nil {
		switch productID {
		case constants.MOD_100_REV1:
			extraCGWMOD100REV1(h, req)
			break
		case constants.FOTX_200_REV1:
			extraCGWFOTX200REV1(h, req)
			break
		case constants.FORX_200_REV1:
			extraCGWFORX200REV1(h, req)
			break
		}
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(404)
}

func UpdateFirmwareHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	serial := vars["serialWorker"][2:]
	var wg sync.WaitGroup
	wg.Add(1)

	file, handler, err := req.FormFile("cyacdFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	extension := filepath.Ext(handler.Filename)
	path := "/data/nms-server/" + handler.Filename

	key := utils.GetEnvVar("bootloader-key")
	if key == "" {
		w.WriteHeader(500)
		w.Write([]byte("ERROR: No key configured.\r\nGo to Configuration to setup the bootloader key."))
		return
	}


	if handler.Size < 1024*1024  && extension == ".cyacd"{
		work := worker.NewSerialWorkRequest(models.Hardware{Serial: serial}, constants.SerialWorkRequestExecCommand, constants.RESET_COMMAND, &wg)
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()

		//usb.ResetCommand(serialWorker)

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer f.Close()
		io.Copy(f, file)


		wg.Add(1)
		go func (w http.ResponseWriter){
			defer wg.Done()
			args := []string{"-path=" + path, "-serialWorker=" + serial, "-key=" + key}
			cmd := exec.Command("bootloader-usb", args...)
			out, e := cmd.CombinedOutput()
			if e !=nil {
				log.Println(e)
				w.WriteHeader(500)
			}else{
				w.WriteHeader(200)
			}
			w.Write(out)
		}(w)
		wg.Wait()
		return
	}
	w.WriteHeader(500)
}

func extraCGWFOTX200REV1(d *models.Hardware, req *http.Request){
	nameExtra := req.PostFormValue("nameExtra")
	command := ""

	var wg sync.WaitGroup

	switch nameExtra {
	case "Laser APC":
		wg.Add(4)
		var laserAPCVal string

		if laserAPCVal = req.PostFormValue("LaserAPC"); laserAPCVal == "" {
			laserAPCVal = "OFF"
		}

		command = fmt.Sprintf("%s %s", constants.SET_LASER_APC_COMMAND, laserAPCVal)
		work := worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)

		command = fmt.Sprintf("%s %s", constants.WRITE_DAC_LASER_BIAS_COMMAND, req.PostFormValue("LaserBiasDAC"))
		work = worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)

		command = fmt.Sprintf("%s %s", constants.WRITE_OPT_DET_OFFSET_COMMAND, req.PostFormValue("OptDetOffset"))
		work = worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)

		command = fmt.Sprintf("%s %s", constants.WRITE_OPT_DET_TARGET_COMMAND, req.PostFormValue("OptDetTarget"))
		work = worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)

		wg.Wait()
		break

	case "Optical Power":
		wg.Add(1)
		//writeCalOptPwr: <wavelength> <intercept> <slope> <AlarmHi> <AlarmLo>
		command = fmt.Sprintf("%s %s %s %s %s %s", constants.WRITE_CAL_OPT_PWR_COMMAND,
			req.PostFormValue("Wavelength"),
			req.PostFormValue("Intercept"),
			req.PostFormValue("Slope"),
			req.PostFormValue("AlarmThresholdHi"),
			req.PostFormValue("AlarmThresholdLo"),
			)
		work := worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
		serialWorker.GetQueue().AddWork(work)
		wg.Wait()
		break
	}
}

func extraCGWFORX200REV1(d *models.Hardware, req *http.Request){
	nameExtra := req.PostFormValue("nameExtra")
	command := ""
	var wg sync.WaitGroup
	wg.Add(1)
	switch nameExtra {
	case "Optical Power":
		//writeCalOptPwr: <wavelength> <intercept> <slope> <AlarmHi> <AlarmLo>
		command = fmt.Sprintf("%s %s %s %s %s %s", constants.WRITE_CAL_OPT_PWR_COMMAND,
			req.PostFormValue("Wavelength"),
			req.PostFormValue("Intercept"),
			req.PostFormValue("Slope"),
			req.PostFormValue("AlarmThresholdHi"),
			req.PostFormValue("AlarmThresholdLo"),
		)
		break
	}
	work := worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
	serialWorker.GetQueue().AddWork(work)
	wg.Wait()
}

func extraCGWMOD100REV1(d *models.Hardware, req *http.Request){
	nameExtra := req.PostFormValue("nameExtra")
	command := ""
	var wg sync.WaitGroup
	wg.Add(1)

	switch nameExtra {
	case "RF AGC":
		var AGCModeVal string

		if AGCModeVal = req.PostFormValue("AgcMode"); AGCModeVal == "" {
			AGCModeVal = "OFF"
		}
		//writeCalRFDet: <frequency> <intercept> <slope> <AlarmHi> <AlarmLo>
		command = fmt.Sprintf("%s %s %s %s %s %s %s",
			constants.WRITE_CAL_RF_AGC_COMMAND,
			AGCModeVal,
			req.PostFormValue("AgcMaxGain"),
			req.PostFormValue("AgcTarget"),
			req.PostFormValue("SquelchMode"),
			req.PostFormValue("SquelchThresholdHi"),
			req.PostFormValue("SquelchThresholdLo"),
		)
		break

	case "RF Detector 2":
		//writeCalRFDet2: <frequency> <intercept> <slope> <AlarmHi> <AlarmLo>
		command = fmt.Sprintf("%s %s %s %s %s %s",
			constants.WRITE_CAL_RF_DET2_COMMAND,
			req.PostFormValue("Frequency"),
			req.PostFormValue("Intercept"),
			req.PostFormValue("Slope"),
			req.PostFormValue("AlarmThresholdHi"),
			req.PostFormValue("AlarmThresholdLo"),
		)
		break
	}
	work := worker.NewSerialWorkRequest(*d, constants.SerialWorkRequestExecCommand, command + "\n", &wg)
	serialWorker.GetQueue().AddWork(work)
	wg.Wait()
}