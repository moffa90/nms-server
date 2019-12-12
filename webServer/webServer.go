package webServer

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/moffa90/triadNMS/assets"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/utils/security"
	"github.com/moffa90/triadNMS/webServer/configuration"
	"github.com/moffa90/triadNMS/webServer/hardware"
	"github.com/moffa90/triadNMS/webServer/home"
	"github.com/moffa90/triadNMS/webServer/information"
	"github.com/moffa90/triadNMS/webServer/login"
	"github.com/moffa90/triadNMS/webServer/logout"
	"github.com/moffa90/triadNMS/webServer/networks/ethernet"
	"github.com/moffa90/triadNMS/webServer/remotes"
	"github.com/moffa90/triadNMS/webServer/snmp"
	"github.com/moffa90/triadNMS/webServer/users"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func LaunchServer(wait time.Duration) {
	security.InitCookieStore()

	demo := assetfs.AssetFS{
		Asset:     assets.Asset,
		AssetDir:  assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix:    "goAssets/demo",
	}

	assets := assetfs.AssetFS{
		Asset:     assets.Asset,
		AssetDir:  assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix:    "goAssets/assets",
	}

	r := mux.NewRouter().StrictSlash(true)
	//Routes
	r.HandleFunc("/login", login.Handler).Methods("GET")
	r.HandleFunc("/auth", login.Authenticate).Methods("POST")

	r.Handle("/", security.CookieMiddleware(http.HandlerFunc(home.Handler))).Methods("GET")
	r.HandleFunc("/403", Handler403).Methods("GET")
	r.Handle("/logout", security.CookieMiddleware(http.HandlerFunc(logout.Handler))).Methods("GET")
	r.Handle("/information", security.CookieMiddleware(http.HandlerFunc(information.Handler))).Methods("GET")
	r.Handle("/configuration", security.CookieMiddleware(http.HandlerFunc(configuration.Handler))).Methods("GET")
	r.Handle("/configuration", security.CookieMiddleware(http.HandlerFunc(configuration.HandlerSave))).Methods("POST")

	//Management routes
	manageRouter := r.PathPrefix("/management").Subrouter().StrictSlash(true)
	manageRouter.Use(security.CookieMiddleware,)

	//Hardware routes
	hardwareRouter := manageRouter.PathPrefix("/hardware").Subrouter().StrictSlash(true)
	hardwareRouter.Use(security.CookieMiddleware)
	hardwareRouter.HandleFunc("/", hardware.Handler).Methods("GET")
	hardwareRouter.HandleFunc("/", hardware.SaveHandler).Methods("POST")
	hardwareRouter.HandleFunc("/{serialWorker}", hardware.InfoHandler).Methods("GET")
	hardwareRouter.HandleFunc("/{serialWorker}", hardware.EditValuesHandler).Methods("POST")
	hardwareRouter.HandleFunc("/{serialWorker}", hardware.DeleteHandler).Methods("DELETE")
	hardwareRouter.HandleFunc("/{serialWorker}/{action}", hardware.ActionHandler).Methods("GET")
	hardwareRouter.HandleFunc("/{serialWorker}/thresholds", hardware.ThresholdsHandler).Methods("POST")
	hardwareRouter.HandleFunc("/{serialWorker}/rfdetector", hardware.RfDetectorHandler).Methods("POST")
	hardwareRouter.HandleFunc("/{serialWorker}/extra", hardware.ExtraHandler).Methods("POST")
	hardwareRouter.HandleFunc("/{serialWorker}/updateFirmware", hardware.UpdateFirmwareHandler).Methods("POST")

	//user routes
	userRouter := manageRouter.PathPrefix("/users").Subrouter().StrictSlash(true)
	userRouter.Use(security.CookieMiddleware, security.AdminMiddleware)
	userRouter.HandleFunc("/", users.Handler).Methods("GET")
	userRouter.HandleFunc("/{userId}", users.SaveHandler).Methods("POST")
	userRouter.HandleFunc("/block/{userId}", users.BlockHandler).Methods("GET")
	userRouter.HandleFunc("/edit/{userId}",  users.EditHandler).Methods("GET")

	//SNMPWorker routes
	snmpRouter := manageRouter.PathPrefix("/snmp").Subrouter().StrictSlash(true)
	snmpRouter.Use(security.CookieMiddleware, security.AdminMiddleware)
	snmpRouter.HandleFunc("/", snmp.Handler).Methods("GET")
	snmpRouter.HandleFunc("", snmp.SaveHandler).Methods("POST")

	//Network routes
	networkRouter := r.PathPrefix("/networks").Subrouter().StrictSlash(true)
	networkRouter.Use(security.CookieMiddleware)

	//WIFI
	//networkRouter.HandleFunc("/wifi", wifi.Handler).Methods("GET")
	//networkRouter.HandleFunc("/wifi", wifi.SaveConfHandler).Methods("POST")
	//networkRouter.HandleFunc("/wifi/success", wifi.SuccessHandler).Methods("GET")
	//networkRouter.HandleFunc("/wifi/error", wifi.ErrorHandler).Methods("GET")

	//ETHERNET
	networkRouter.HandleFunc("/ethernet", ethernet.Handler).Methods("GET")
	networkRouter.HandleFunc("/ethernet", ethernet.SaveConfHandler).Methods("POST")
	networkRouter.HandleFunc("/ethernet/success", ethernet.SuccessHandler).Methods("GET")
	networkRouter.HandleFunc("/ethernet/error", ethernet.ErrorHandler).Methods("GET")

	mode := os.Getenv("mode")
	if strings.HasPrefix(mode, "hec-") {
		//Remotes Routes
		remotesRouter := manageRouter.PathPrefix("/remotes").Subrouter().StrictSlash(true)
		remotesRouter.Use(security.CookieMiddleware)
		remotesRouter.HandleFunc("/", remotes.Handler).Methods("GET")
		remotesRouter.HandleFunc("/", remotes.AddRemoteHandler).Methods("POST")
		remotesRouter.HandleFunc("/{remote}/{group}", remotes.DeleteHandler).Methods("DELETE")
	}


	//Resources
	r.PathPrefix("/demo/").Handler(http.StripPrefix("/demo/", http.FileServer(&demo))).Methods("GET")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(&assets))).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	db.Shared.Close()
	os.Exit(0)
}

func Handler403(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]bool)
	forbidden, _ := assets.Asset("goAssets/html/403.html")
	tpl := template.New("templates")
	tpl.New("403").Parse(string(forbidden))

	if ErrTpl :=  tpl.ExecuteTemplate(w, "403", data); ErrTpl != nil{
		log.Println(ErrTpl.Error())
	}
}