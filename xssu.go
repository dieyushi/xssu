package main

import (
	"code.google.com/p/gorilla/sessions"
	"database/sql"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nranchev/go-libGeoIP"
	"io"
	"log"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte("something-very-secret-for-session"))
var orm beedb.Model
var gi *libgeo.GeoIP

func main() {
	handleConfig()

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			LogE(err)
		}
		log.SetOutput(io.MultiWriter(f, os.Stderr))
		defer f.Close()
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		LogE(err)
	}
	orm = beedb.New(db)

	gi, err = libgeo.Load(geoipdb)
	if err != nil {
		gi, err = libgeo.Load("/usr/share/GeoIP/GeoLiteCity.dat")
		if err != nil {
			Log(err)
			return
		}
	}

	http.Handle("/js/", noDirListing(http.FileServer(http.Dir("templates"))))
	http.Handle("/css/", noDirListing(http.FileServer(http.Dir("templates"))))
	http.Handle("/pic/", noDirListing(http.FileServer(http.Dir("templates"))))
	http.Handle("/img/", noDirListing(http.FileServer(http.Dir("templates"))))

	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/about/", aboutHandler)
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/reset/password/", resetpasswordHandler)
	http.HandleFunc("/reset/email/", resetemailHandler)
	http.HandleFunc("/setting/", settingHandler)
	http.HandleFunc("/addworker/", addworkerHandler)
	http.HandleFunc("/editworker/", editworkerHandler)
	http.HandleFunc("/delworker/", delworkerHandler)
	http.HandleFunc("/delrecord/", delrecordHandler)
	http.HandleFunc("/worker/", workerHandler)
	http.HandleFunc("/modules/", modulesHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/u/", uHandler)
	http.HandleFunc("/s/", sHandler)
	http.HandleFunc("/", notfoundHandler)

	host = os.Args[1]
	err = http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		LogE(err)
	}
}
