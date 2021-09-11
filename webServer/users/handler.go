package users

import (
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/utils"
	"github.com/moffa90/triadNMS/utils/security"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	infoStruct := struct{
		Users []models.User
		Active string
		Info map[string]string
	}{
		models.GetUsers(db.Shared),
		"management",
		security.CookieGetInfo(w, req),
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_USERS_PATH, infoStruct)
}

func BlockHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["userId"]
	if Err :=  models.ChangeStatusUser(db.Shared, userId); Err != nil{
		log.Println(Err.Error())
	}
	http.Redirect(w, req, "/management/users", http.StatusFound)
}

func EditHandler(w http.ResponseWriter, req *http.Request) {
	//TODO: finish edit user
	vars := mux.Vars(req)
	userId := vars["userId"]

	if usr, err := models.GetUserById(db.Shared, userId); err == nil {
		infoStruct := struct {
			Roles  []models.Role
			User   models.User
			Active string
			Info   map[string]string
		}{
			models.GetRoles(db.Shared),
			usr,
			"management",
			security.CookieGetInfo(w, req),
		}

		utils.RenderPage(w, constants.TEMPLATE_PAGE_EDIT_USERS_PATH, infoStruct)

	}else {
		http.Redirect(w, req, "/users", http.StatusFound)
	}
}

func SaveHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["userId"]
	name := req.PostFormValue("name")
	password := req.PostFormValue("password")
	password_confirm := req.PostFormValue("password-confirm")
	role := req.PostFormValue("role")

	if name == "" {
		req.PostForm.Add("error", "Empty Name")
		http.Redirect(w, req, "/users/edit/"+userId, http.StatusTemporaryRedirect)
		return
	}

	if password != "" && password_confirm != ""{
		if password != password_confirm{
			req.PostForm.Add("error", "Passwords are not equal")
			http.Redirect(w, req, "/users/edit/"+userId, http.StatusTemporaryRedirect)
			return
		}
	}

	if role == ""{
		req.PostForm.Add("error", "Empty Role")
		http.Redirect(w, req, "/users/edit/"+userId, http.StatusTemporaryRedirect)
		return
	}



	w.Write([]byte(userId + name + password + password_confirm + role))
}