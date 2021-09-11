package wifi

import (
	"github.com/moffa90/triadNMS/assets"
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/utils"
	"github.com/moffa90/triadNMS/utils/security"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	if i, err := models.GetInterfaceByName(db.Shared, "wlan0"); err != nil{
		http.Redirect(w, req, "/", http.StatusFound)
	}else{
		infoStruct := struct{
			Active string
			Info map[string]string
			InterfaceInfo models.Interface
			Errors map[string]string
		}{
			"networks",
			security.CookieGetInfo(w, req),
			i,
			nil,
		}

		utils.RenderPage(w, constants.TEMPLATE_PAGE_WIFI_PATH, infoStruct)
	}


}

func SaveConfHandler(w http.ResponseWriter, req *http.Request) {
	if i, err := models.GetInterfaceByName(db.Shared, "wlan0"); err != nil{
		http.Redirect(w, req, "/", http.StatusFound)
	}else{
		errors := make(map[string]string)

		ssid := req.PostFormValue("ssid")
		password := req.PostFormValue("password")
		active := req.PostFormValue("active")

		if active, err := strconv.ParseBool(active); err == nil {
			i.DHCP = active
			if active{
				if ssid == "" {
					errors["1ssid"] = "SSID field could not be empty"
				}

				if password == ""{
					errors["2pass"] = "Password field could not be empty"
				}
				i.SSID = ssid
				i.Password = password
			}

			if len(errors) == 0 {
				if active{
					//activate ap

					tpl := template.New("temp")
					wpa, _ := assets.Asset(constants.TEMPLATE_AP_PATH)
					utils.BackupFile(constants.LOCAL_AP_CONF_PATH)

					f, errFile := utils.CreateFile(constants.LOCAL_AP_CONF_PATH)

					if errFile != nil {
						log.Println("Error creating File:" + errFile.Error())
					}

					tpl.New("ap").Parse(string(wpa))

					if ErrTpl := tpl.ExecuteTemplate(f, "ap", i); ErrTpl != nil {
						log.Println(ErrTpl.Error())
					}

					f.Close()
					_, errScript := utils.ExecScript(constants.START_AP_SCRIPT_PATH)
					if errScript != nil {
						log.Println(errScript.Error())
						errors["errScript"] = fmt.Sprintf("AP script error: %s\n", errScript.Error())
					}else{
						db.Shared.Save(&i)
						http.Redirect(w, req, "wifi/success?msg=active", http.StatusFound)
						return
					}
				}else {
					_, errScript := utils.ExecScript(constants.STOP_AP_SCRIPT_PATH)
					if errScript != nil {
						log.Println(errScript.Error())
						errors["errScript"] = fmt.Sprintf("AP script error: %s\n", errScript.Error())
					}else{
						db.Shared.Save(&i)
						http.Redirect(w, req, "wifi/success?msg=deactive", http.StatusFound)
						return
					}
				}

			}
		}



		infoStruct := struct{
			Active string
			Info map[string]string
			InterfaceInfo models.Interface
			Errors map[string]string
		}{
			"networks",
			security.CookieGetInfo(w, req),
			i,
			errors,
		}


		utils.RenderPage(w, constants.TEMPLATE_PAGE_WIFI_PATH, infoStruct)
		return
	}

}

func ErrorHandler(w http.ResponseWriter, req *http.Request)  {
	msg := req.FormValue("msg")

	switch (msg){
	case "failed":
		msg = "The AP activation failed"
		break
	}

	infoStruct := struct{
		Active string
		Info map[string]string
		Message string
	}{
		"networks",
		security.CookieGetInfo(w, req),
		msg,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_ERROR_PATH, infoStruct)
}

func SuccessHandler(w http.ResponseWriter, req *http.Request)  {

	msg := req.FormValue("msg")

	switch (msg){
	case "active":
		msg = "The AP was activated successfully"
		break
	case "deactive":
		msg = "The AP was deactivated successfully"
		break
	}

	infoStruct := struct{
		Active string
		Info map[string]string
		Message string
	}{
		"networks",
		security.CookieGetInfo(w, req),
		msg,
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_SUCCESS_PATH, infoStruct)
}