package user

import "github.com/dgrijalva/jwt-go"

var JwtKey = []byte("my_secret_key")

type Claims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}
