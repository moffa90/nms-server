package remotes

import (
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"github.com/moffa90/nms-server/utils/snmpClient"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	infoStruct := struct{
		Active string
		Info map[string]string
		Remotes [][]models.Remote
	}{
		"management",
		security.CookieGetInfo(w, req),
		nil,
	}

	remotes := models.GetRemotes(db.Shared)
	infoStruct.Remotes = make([][]models.Remote, 4)

	for i, _ := range infoStruct.Remotes{
		infoStruct.Remotes[i] = make([]models.Remote, 5)
	}

	for _, r := range remotes{
		infoStruct.Remotes[r.Remote][r.Group] = r
	}
	utils.RenderPage(w, constants.TEMPLATE_PAGE_REMOTES_PATH, infoStruct)
}

func AddRemoteHandler(w http.ResponseWriter, req *http.Request){
	ip := req.PostFormValue("ip-address")
	port := req.PostFormValue("port")
	group := req.PostFormValue("group")
	remote := req.PostFormValue("remote")
	errors := make(map[string]string)

	addressIP := net.ParseIP(ip)
	if addressIP.To4() != nil {
		remoteInt, _ := strconv.ParseInt(remote, 10, 8)
		groupInt, _ := strconv.ParseInt(group, 10, 8)
		r := models.Remote{
			Remote:   int8(remoteInt),
			Group:    int8(groupInt),
			Ip:       ip,
			Port:     port,
			Hostname: "",
		}

		if hostname, e := snmpClient.GetHostnameRemoteSNMP(ip, port, "public"); e != nil{
			errors["0"] += e.Error()

		}else{
			r.Hostname = hostname
			if e := models.CreateOrUpdate(db.Shared, &r); e != nil{
				errors["0"] += e.Error()
			}
		}

	}else{
		errors["1"] = fmt.Sprintf("IP Address: \"%v\" is not an IPv4 address\n", ip)
	}

	if len(errors) > 0{
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("status","500")
		if payload, err := json.Marshal(errors); err == nil {
			http.Error(w, string(payload), http.StatusBadRequest)
		}
	}else{
		w.WriteHeader(200)
	}


}

func DeleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	group :=  vars["group"]
	remote :=  vars["remote"]

	if err := db.Shared.Debug().Unscoped().Delete(models.Remote{}, "remote == ? AND `group` == ?", remote, group).Error; err == nil{
		w.WriteHeader(200)
	}else{
		log.Println(err.Error())
	}

}