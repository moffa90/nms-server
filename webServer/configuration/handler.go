package configuration

import (
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"encoding/hex"
	"log"
	"net/http"
	"os/exec"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	c,_ := models.GetConfigByCat("general", db.Shared)
	for i := 0; i < len(c); i ++ {
		if c[i].Key == "bootloader-key" {
			c[i].Value = utils.GetEnvVar("bootloader-key")
		}
	}

	infoStruct := struct{
		Active string
		Info map[string]string
		Configurations []models.Configuration
	}{
		"config",
		security.CookieGetInfo(w, req),
		c,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_CONFIG_PATH, infoStruct)
}

func HandlerSave(w http.ResponseWriter, req *http.Request) {
	hostname := req.PostFormValue("hostname")
	bootloaderKey := req.PostFormValue("bootloader-key")

	if conf, err := models.GetConfigByKey("hostname", db.Shared); err == nil {
		if hostname != "" {
			conf.Value = hostname
			if err = db.Shared.Save(&conf).Error; err != nil {
				log.Println(err.Error())
			}else{
				cmd := exec.Command("hostnamectl", "set-hostname", hostname)
				e := cmd.Run()
				if e != nil {
					log.Println(e.Error())
				}
			}
		}
	}else{
		log.Println(err.Error())
	}

	if _, e := hex.DecodeString(bootloaderKey); e == nil && len(bootloaderKey) == 12{
		utils.SetEnvVar("bootloader-key", bootloaderKey)
	}

	c,_ := models.GetConfigByCat("general", db.Shared)
	for i := 0; i < len(c); i ++ {
		if c[i].Key == "bootloader-key" {
			c[i].Value = utils.GetEnvVar("bootloader-key")
		}
	}

	infoStruct := struct{
		Active string
		Info map[string]string
		Configurations []models.Configuration
	}{
		"config",
		security.CookieGetInfo(w, req),
		c,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_CONFIG_PATH, infoStruct)
}

