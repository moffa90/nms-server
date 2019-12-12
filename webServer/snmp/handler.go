package snmp

import (
	"github.com/moffa90/nms-server/assets"
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"text/template"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	c, _ := models.GetConfigByCat("snmp", db.Shared)
	infoStruct := struct{
		Active string
		Info map[string]string
		Errors map[string]string
		SNMP []models.Configuration
	}{
		"snmp",
		security.CookieGetInfo(w, req),
		nil,
		c,
		}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_SNMP_PATH, infoStruct)
}

func SaveHandler(w http.ResponseWriter, req *http.Request) {

	var port uint64 = 0
	var e error

	portStr := req.PostFormValue("port")
	writeString := req.PostFormValue("write-string")
	readString := req.PostFormValue("read-string")
	sysContact := req.PostFormValue("sys-contact")
	sysLocation := req.PostFormValue("sys-location")
	errors := make(map[string]string)

	if port, e = strconv.ParseUint(portStr, 10, 64); e == nil {
		if port > 0xffff {
			errors["port"] = "Port must be less then 65535"
		}
	}

	if writeString == "" {
		errors["writeString"] = "Write string cannot be empty"
	}

	if readString == "" {
		errors["readString"] = "Read string cannot be empty"
	}

	if sysContact == "" {
		errors["syscontact"] = "System Contact string cannot be empty"
	}

	if sysLocation == "" {
		errors["syslocation"] = "System Location string cannot be empty"
	}

	infoStruct := struct{
		Active string
		Info map[string]string
		Errors map[string]string
		SNMP []models.Configuration
	}{
		"snmp",
		security.CookieGetInfo(w, req),
		errors,
		nil,
	}

	snmpdInfo := struct {
		WriteString string
		ReadString string
		Port uint64
		SysContact string
		SysLocation string
	}{
		writeString,
		readString,
		port,
		sysContact,
		sysLocation,
	}

	if len(errors) == 0  {
		if conf, err := models.GetConfigByKey("port", db.Shared); err == nil {
			conf.Value = portStr
			db.Shared.Save(&conf)
		}

		if conf, err := models.GetConfigByKey("write-string", db.Shared); err == nil {
			conf.Value = writeString
			db.Shared.Save(&conf)
		}

		if conf, err := models.GetConfigByKey("read-string", db.Shared); err == nil {
			conf.Value = readString
			db.Shared.Save(&conf)
		}

		if conf, err := models.GetConfigByKey("sys-contact", db.Shared); err == nil {
			conf.Value = sysContact
			db.Shared.Save(&conf)
		}

		if conf, err := models.GetConfigByKey("sys-location", db.Shared); err == nil {
			conf.Value = sysLocation
			db.Shared.Save(&conf)
		}

		c, _ := models.GetConfigByCat("snmp", db.Shared)
		infoStruct.SNMP = c

		snmpd, _ := assets.Asset(constants.TEMPLATE_SNMPD_CONF_PATH)
		utils.BackupFile(constants.LOCAL_SNMPD_CONF_PATH)

		f, errFile := utils.CreateFile(constants.LOCAL_SNMPD_CONF_PATH)

		if errFile != nil {
			log.Println("Error creating File:" + errFile.Error())
		}

		tpl := template.New("temp")
		tpl.New("snmpd").Parse(string(snmpd))

		if ErrTpl := tpl.ExecuteTemplate(f, "snmpd", snmpdInfo); ErrTpl != nil {
			log.Println(ErrTpl.Error())
		}

		f.Close()

		cmd := exec.Command("systemctl", "restart", "snmpd")
		out, e := cmd.CombinedOutput()
		if e !=nil {
			log.Println(out)
			log.Println(e.Error())
		}

	} else {
		aux := make([]models.Configuration,0)
		aux = append(aux, models.Configuration{Key: "port", Value: portStr})
		aux = append(aux, models.Configuration{Key: "write-string", Value: writeString})
		aux = append(aux, models.Configuration{Key: "read-string", Value: readString})
		aux = append(aux, models.Configuration{Key: "sys-contact", Value: sysContact})
		aux = append(aux, models.Configuration{Key: "sys-location", Value: sysLocation})

		infoStruct.SNMP = aux
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_SNMP_PATH, infoStruct)
}
