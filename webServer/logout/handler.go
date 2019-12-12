package logout

import (
	"github.com/moffa90/nms-server/utils/security"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	sess, err := security.Store.Get(req, "session")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	sess.Options.MaxAge = -1

	err = sess.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, req, "/login", http.StatusFound)
}