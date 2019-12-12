package ethernet

import (
	"github.com/moffa90/nms-server/assets"
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
)

//TODO: On rollback, put previous configuration instead DHCP... ask

func setDHCP(){
	tpl := template.New("temp")
	eth, _ := assets.Asset(constants.TEMPLATE_ETH_DHCP_PATH)
	tpl.New("eth-dhcp").Parse(string(eth))

	utils.BackupFile(constants.LOCAL_ETH_CONF_PATH)

	f, errFile := utils.CreateFile(constants.LOCAL_ETH_CONF_PATH)

	if errFile != nil {
		log.Println("Error creating File:" + errFile.Error())
	}

	if ErrTpl :=  tpl.ExecuteTemplate(f, "eth-dhcp", nil); ErrTpl != nil{
		log.Println(ErrTpl.Error())
	}

	f.Close()
}

func setStatic(i models.Interface){
	tpl := template.New("temp")
	info := struct {
		IPAddress   string
		Gateway    string
		SubnetMask string
	}{
		i.IPAddress,
		i.Gateway,
		i.SubnetMask,
	}
	eth, _ := assets.Asset(constants.TEMPLATE_ETH_STATIC_PATH)

	utils.BackupFile(constants.LOCAL_ETH_CONF_PATH)

	f, errFile := utils.CreateFile(constants.LOCAL_ETH_CONF_PATH)

	if errFile != nil {
		log.Println("Error creating File:" + errFile.Error())
	}

	tpl.New("eth-static").Parse(string(eth))

	if ErrTpl := tpl.ExecuteTemplate(f, "eth-static", info); ErrTpl != nil {
		log.Println(ErrTpl.Error())
	}

	f.Close()
}

func Handler(w http.ResponseWriter, req *http.Request) {

	if i, err := models.GetInterfaceByName(db.Shared, "eth0"); err != nil{
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

		utils.RenderPage(w, constants.TEMPLATE_PAGE_ETHERNET_PATH, infoStruct)
	}

}

func SaveConfHandler(w http.ResponseWriter, req *http.Request) {
	errors := make(map[string]string)

	if i, err := models.GetInterfaceByName(db.Shared, "eth0"); err == nil {
		//get data from http form
		dhcp := req.PostFormValue("dhcp")

		if dhcp, err := strconv.ParseBool(dhcp); err == nil {
			i.DHCP = dhcp
			if !dhcp{
				address := req.PostFormValue("ip-address")
				subnet := req.PostFormValue("subnet-mask")
				gateway := req.PostFormValue("gateway")

				addressIP := net.ParseIP(address)
				if addressIP.To4() == nil {
					errors["0"] = fmt.Sprintf("IP Address: \"%v\" is not an IPv4 address\n", address)
				}

				subnetIP := net.ParseIP(subnet)
				if subnetIP.To4() == nil {
					errors["1"] = fmt.Sprintf("Subnet Mask: \"%v\" is not an IPv4 address\n", subnet)
				}

				gatewayIP := net.ParseIP(gateway)
				if gatewayIP.To4() == nil {
					errors["2"] = fmt.Sprintf("Gateway: \"%v\" is not an IPv4 address\n", gateway)
				}


				i.IPAddress = address
				i.Gateway = gateway
				i.SubnetMask = subnet
				if len(errors) == 0 {
					setStatic(i)
				}
			} else {
				setDHCP()
			}

			if len(errors) == 0{
				_, errScript := utils.ExecScript(constants.ETH_SCRIPT_PATH)
				if errScript != nil {
					log.Println(errScript.Error())
					errors["errScript"] = fmt.Sprintf("Eth0 ping script error: %s\n", errScript.Error())
				}else{
					db.Shared.Save(&i)
					http.Redirect(w, req, "ethernet/success", http.StatusFound)
					return
					/*} else {
						setDHCP()
						utils.ExecScript(constants.ETH_SCRIPT_PATH)
						i.DHCP = true
						db.Shared.Save(&i)
						http.Redirect(w, req, "ethernet/error?msg=staticFailed", http.StatusFound)
						return
					}*/
				}
			}
		}else{
			log.Println(err.Error())
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

		utils.RenderPage(w, constants.TEMPLATE_PAGE_ETHERNET_PATH, infoStruct)
		return
	}

}

func ErrorHandler(w http.ResponseWriter, req *http.Request)  {
	msg := req.FormValue("msg")

	if msg == "staticFailed"{
		msg = "The custom configuration failed, the device was set back to DHCP"
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
	infoStruct := struct{
		Active string
		Info map[string]string
		Message string
	}{
		"networks",
		security.CookieGetInfo(w, req),
		"The ethernet interface was configured successfully",
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_SUCCESS_PATH, infoStruct)
}