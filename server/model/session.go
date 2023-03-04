package model

import (
	"github.com/jackc/pgtype"
	"html/template"
	"net"
	"net/http"
	"sdu.store/server"
	"time"
)

type Session struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	UUID      string      `json:"uuid"`
	CreatedAt time.Time   `json:"created_at"`
	DeletedAt time.Time   `json:"deleted_at"`
	LastLogin time.Time   `json:"last_login"`
	IP        pgtype.Inet `json:"ip"`
}

var CurrentSession = make(map[string]Session)

func CreateSession() {

}

func ResolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {

			ip := networkIp.IP.String()

			return ip
		}
	}
	return ""
}

func SetInet() pgtype.Inet {
	var inet pgtype.Inet
	if err := inet.Set(ResolveHostIp()); err != nil {
		panic(err)
	}
	return inet
}

func GetAllSessions(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/AdminSession.gohtml")
	var sessions []Session
	server.DB.Find(&sessions)
	if err := tm.Execute(w, sessions); err != nil {
		panic(err)
	}
}
