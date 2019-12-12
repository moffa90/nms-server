package login

import (
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]bool)

	if req.URL.Query().Get("failedLogin") == "true" {
		data["failedLogin"] = true
	}

	utils.RenderPage(w , constants.TEMPLATE_PAGE_LOGIN_PATH, data)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	password := req.FormValue("password")

	userObj, err := models.GetUserByUsername(db.Shared, user)
	if err != nil{
		req.Form.Add("failedLogin", "true")
		http.Redirect(w, req, "/login?failedLogin=true", http.StatusFound)
		return
	}

	if !userObj.Active{
		http.Redirect(w, req, "/login?failedLogin=true", http.StatusFound)
		return
	}

	if security.ComparePasswords(userObj.Password, []byte(password)) {
		session, _ := security.Store.Get(req, "session")
		session.Values["authenticated"] = true
		session.Values["role"] = userObj.Role.Name
		session.Values["userId"] = userObj.Id
		session.Values["name"] = userObj.Name
		session.Values["username"] = userObj.Username
		session.Save(req, w)

		userObj.LastLogin = time.Now()
		db.Shared.Save(&userObj)
		http.Redirect(w, req, "/", http.StatusFound)
		return
	} else {
		req.Form.Add("failedLogin", "true")
		http.Redirect(w, req, "/login?failedLogin=true", http.StatusFound)
		return
	}
}