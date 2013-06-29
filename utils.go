package main

import (
	"crypto/md5"
	"crypto/sha512"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
		if strings.HasSuffix(req.URL.Path, "/") {
			http.NotFound(rw, req)
			return
		}
		h.ServeHTTP(rw, req)
	})
}

func CheckSession(req *http.Request) bool {
	session, _ := store.Get(req, "session")
	logged := session.Values["logged"]
	if logged == nil || logged.(string) != "1" {
		return false
	}
	return true
}

func CheckPassword(req *http.Request, password string) bool {
	var u XssuUser
	orm.Where("userid = ?", GetUserID(req)).Find(&u)

	if u.Password == HashPassword(password, u.Password[:16]) {
		return true
	}
	return false
}

func CheckHashID(hashid string) bool {
	var worker XssuWorker
	orm.Where("hashid = ?", hashid).Find(&worker)
	if worker.Userid != 0 {
		return true
	}
	return false
}

func GetUserID(req *http.Request) int {
	session, _ := store.Get(req, "session")
	userid := session.Values["userid"]
	return userid.(int)
}

func GenSalt(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	randomPassword := make([]byte, length)

	for i := 0; i < len(randomPassword); i++ {
		randomPassword[i] = alphabet[rand.Int()%len(alphabet)]
	}
	return string(randomPassword)
}

func HashID(unixtime int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	h := md5.New()
	io.WriteString(h, strconv.Itoa(unixtime)+GenSalt(8))
	hex := fmt.Sprintf("%x", h.Sum(nil))
	res := make([]string, 4)

	for i := 0; i < 4; i++ {
		hexint, _ := strconv.ParseInt(hex[i*8:(i+1)*8], 16, 64)
		hexint = 0x3FFFFFFF & hexint
		var outChars string

		for j := 0; j < 6; j++ {
			index := 0x0000003D & hexint
			outChars += string(chars[index])
			hexint = hexint >> 5
		}
		res[i] = outChars
	}
	return res[0]
}

func HashPassword(password string, salt string) string {
	h := sha512.New()
	io.WriteString(h, password+salt)
	hex := fmt.Sprintf("%x", h.Sum(nil))
	return salt + hex[:32]
}

func ShowUnixTime(timeUnix int) string {
	t := time.Unix(int64(timeUnix), 0)
	return t.Format("2006-01-02 15:04:05")
}

func GetIPLocation(ip string) string {
	loc := gi.GetLocationByIP(ip)
	if loc != nil {
		return loc.CountryName + "/" + loc.City
	}
	return ""
}
