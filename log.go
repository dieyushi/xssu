package main

import (
	"log"
)

func LogW(v ...interface{}) {
	if weblog {
		log.SetPrefix("[WEB] ")
		log.Println(v...)
	}
}

func Log(v ...interface{}) {
	log.SetPrefix("[XSSU] ")
	log.Println(v...)
}

func LogE(v ...interface{}) {
	log.SetPrefix("[XSSU] ")
	log.Fatalln(v...)
}
