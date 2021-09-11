package main

import (
	"github.com/moffa90/triadNMS/assets"
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/subagentSNMP"
	"github.com/moffa90/triadNMS/utils"
	"github.com/moffa90/triadNMS/utils/usb"
	"github.com/moffa90/triadNMS/webServer"
	"github.com/moffa90/triadNMS/worker/SNMPWorker"
	//"github.com/moffa90/triadNMS/worker/serialWorker"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)




//scp -P 2284 /home/jmoffa/Escritorio/NMS/nms-image/build/tmp/sysroots-components/cortexa7hf-neon-vfpv4/nms-server/usr/bin/nms-server root@cellgain.ddns.net:/home/root

func init() {

	os.MkdirAll(constants.LOCAL_DIRECTORY, os.ModePerm)
	os.MkdirAll(constants.LOCAL_NETWORK_CONF_DIRECTORY, os.ModePerm)

	if !utils.ExistsFile(constants.ENV_FILE_PATH) {
		env, _ := assets.Asset(constants.TEMPLATE_ENV_PATH)
		f,_ := utils.CreateFile(constants.ENV_FILE_PATH)
		f.Write(env)
		f.Close()
	}

	content, _ := ioutil.ReadFile(constants.ENV_FILE_PATH)
	myEnv, err := godotenv.Unmarshal(string(content))
	if err != nil {
		log.Fatal("Error writing .env file: " + err.Error())
	}

	if myEnv["app-key"] == "" {
		key := make([]byte, 64)
		_, err := rand.Read(key)

		if err != nil {
			log.Fatalln("Error creating app-key: " + err.Error())
		}

		strKey := base64.StdEncoding.EncodeToString(key)
		myEnv["app-key"] = strKey
		err = godotenv.Write(myEnv, constants.ENV_FILE_PATH)

		if err != nil {
			log.Fatalln("Error writting .env file: " + err.Error())
		}
	}

	err = godotenv.Load(constants.ENV_FILE_PATH)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Init()
	usb.Init()

	if c, e := models.GetConfigByKey("hostname", db.Shared); e != nil {
		log.Println(e.Error())
	}else{
		cmd := exec.Command("hostnamectl", "set-hostname", c.Value)
		e := cmd.Run()
		if e != nil {
			log.Println(e.Error())
		}
	}
}

//install https://github.com/go-bindata/go-bindata
//go:generate go-bindata -pkg assets -o assets/bindata.go goAssets/...
func main() {
	var wait time.Duration
	modeFlag := flag.String("mode", "hec", "hec-1, hec-2 or remote")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	utils.SetEnvVar("mode", *modeFlag)

	if i, err := models.GetInterfaceByName(db.Shared, "wlan0"); err == nil{
		if i.DHCP{ //In wifi, DHCP field works as active.
			go utils.ExecScript(constants.START_AP_SCRIPT_PATH)
		}
	}


	go subagentSNMP.StartSubagent()
	//go serialWorker.StartWorker()
	go SNMPWorker.StartWorker()

	//var wg sync.WaitGroup
	//wg.Add(1)
	//w := worker.NewSNMPWorkRequest("192.168.0.174", "161", "public", constants.SNMPWorkRequestGetHostname, &wg)
	//wSNMP.GetQueue().AddWork(w)
	//wg.Wait()
	//log.Info(w.Response)

	webServer.LaunchServer(wait)

	//close db before end.

}
