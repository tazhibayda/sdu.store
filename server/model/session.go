package model

import "time"

type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	UUID      int64     `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	LastLogin time.Time `json:"last_login"`
}

func CreateSession() {

}
