package main

type XssuUser struct {
	Username string
	Password string
	Email    string
	Userid   int `PK`
}

type XssuWorker struct {
	Workerid int `PK`
	Userid   int
	Time     int
	Note     string
	Hashid   string
}

type XssuVictims struct {
	Vid    int `PK`
	Hashid string
	Title  string
	Url    string
	Ua     string
	Ip     string
	Time   int
	Cookie string
}

type WorkerInfo struct {
	WorkerID string
	Time     string
	Note     string
	Num      int
}

type HomeData struct {
	WorkerMap    map[string]*WorkerInfo
	CurrentPage  int
	PageNum      int
	NextPage     int
	PreviousPage int
}

type VictimInfo struct {
	VID      string
	Title    string
	Url      string
	Ua       string
	Ip       string
	Time     string
	Cookie   string
	Loc      string
	WorkerID string
}

type VictimData struct {
	VictimMap    map[string]*VictimInfo
	WorkerID     string
	CurrentPage  int
	PageNum      int
	NextPage     int
	PreviousPage int
}
