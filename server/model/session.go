package model

import (
	"github.com/jackc/pgtype"
	"net"
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

func CreateSession() {
}

func (session *Session) Check() bool {
	server.DB.Where("UUID=?", session.UUID).Find(session)
	if session.ID == 0 {
		return false
	}
	if session.DeletedAt.After(time.Now()) {
		session.Delete()
		return false
	}
	return true
}

func (session *Session) CheckStaff() bool {
	server.DB.Where("UUID=?", session.UUID).Find(session)
	if session.ID == 0 {
		return false
	}

	var user User
	server.DB.Where("ID=?", session.UserID).Find(&user)

	return user.Is_staff
}

func (session *Session) Delete() {
	server.DB.Where("ID=?", session.ID).Delete(&session)
}

func (s *Session) Expired() bool {
	return s.DeletedAt.Sub(time.Now()) < 0
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
