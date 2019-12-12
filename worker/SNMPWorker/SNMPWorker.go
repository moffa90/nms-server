package SNMPWorker

import (
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/utils/snmpClient"
	"github.com/moffa90/triadNMS/worker"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

// A buffered channel that we can send work requests on.

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

func StartWorker() {
	go periodicallySNMPWorks()
	GetQueue()
	for {
		work, ok := <-queue.Pending
		if ok {
			// Receive a work request.
			switch work.GetAction() {
			case constants.SNMPWorkRequestGetHostname:
				//go func() {
				start := time.Now()
					var e error
					workAux := work.(*worker.SNMPWorkRequest)
					if workAux.Response, e = snmpClient.GetHostnameRemoteSNMP(work.(*worker.SNMPWorkRequest).Remote.Ip, work.(*worker.SNMPWorkRequest).Remote.Port, work.(*worker.SNMPWorkRequest).Community); e == nil {
						workAux.Remote.Hostname = workAux.Response
						db.Shared.Save(workAux.Remote)
					}

					if work.GetWg() != nil {
						work.GetWg().Done()
					}
				//}()
				log.Printf("End update %d-%d: %s", workAux.Remote.Group, workAux.Remote.Remote, time.Since(start))
				break
			case constants.SNMPWorkRequestGetAlarms:
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

func periodicallySNMPWorks() {
	ticker := time.NewTicker(1 * time.Second)
	var work *worker.SNMPWorkRequest
	for _ = range ticker.C {
		remotes := models.GetRemotes(db.Shared)
		if community, err := models.GetConfigByKey("read-string", db.Shared); err == nil {
			for _, r := range remotes {
				work = worker.NewSNMPWorkRequest(r, community.Value, constants.SNMPWorkRequestGetAlarms, nil)
				queue.AddWork(work)
				work = worker.NewSNMPWorkRequest(r, community.Value, constants.SNMPWorkRequestGetHostname, nil)
				queue.AddWork(work)
			}
		}
	}
}
