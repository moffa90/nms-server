package security

import (
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

var(
	key string
	Store *sessions.CookieStore
)

func InitCookieStore(){
	key = os.Getenv("app-key")
	Store = sessions.NewCookieStore([]byte(key))
}
func HashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "session")
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := CookieGetInfo(w, r)
		// Check if user is authenticated
		if info["role"] == "admin"{
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/403", http.StatusFound)
		}
	})
}

func CookieGetInfo(w http.ResponseWriter, req *http.Request) map[string]string{
	session, err := Store.Get(req, "session")
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, req, "/login", http.StatusFound)
	}

	mapReturn := make(map[string]string)
	mapReturn["role"] = session.Values["role"].(string)
	mapReturn["username"] = session.Values["username"].(string)
	mapReturn["name"] = session.Values["name"].(string)

	return mapReturn
}


