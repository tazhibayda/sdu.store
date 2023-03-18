package model

import (
	"time"
)

type Userdata struct {
	UserId      int64     `json:"user_id"`
	Firstname   string    `json:"first_name"`
	Lastname    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	CountryCode string    `json:"country_code"`
	ZIPCode     string    `json:"zip_code"`
	Birthday    time.Time `json:"birthday"`
}
