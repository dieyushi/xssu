package main

import (
	"code.google.com/p/goconf/conf"
)

var (
	host    string
	domain  string
	port    string
	dbfile  string
	weblog  bool
	logfile string
	geoipdb string
)

func handleConfig() {
	c, err := conf.ReadConfigFile("xssu.ini")
	if err != nil {
		Log("parse config error (xssu.conf not found), start up with default config")
		port = "80"
		dbfile = "xssu.db"
		weblog = true
		logfile = "xssu.log"
		geoipdb = "GeoLiteCity.dat"
		return
	}

    host, err = c.GetString("default", "host")
    if err != nil {
        host = "127.0.0.1"
    }
	port, err = c.GetString("default", "port")
	if err != nil {
		port = "80"
	}
	dbfile, err = c.GetString("default", "dbfile")
	if err != nil {
		dbfile = "xssu.db"
	}
	weblog, err = c.GetBool("default", "weblog")
	if err != nil {
		weblog = true
	}
	logfile, err = c.GetString("default", "logfile")
	if err != nil {
		logfile = "xssu.log"
	}
	geoipdb, err = c.GetString("default", "geoipdb")
	if err != nil {
		geoipdb = "GeoLiteCity.dat"
	}
}
