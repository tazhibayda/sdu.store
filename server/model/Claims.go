package model

import (
	"github.com/golang-jwt/jwt/v5"
	"sdu.store/server"
)

type Claims struct {
	User *User
	jwt.RegisteredClaims
}

// add value to dbINFO
var JwtKey = []byte(server.TokenKey)
