package user

import "github.com/dgrijalva/jwt-go"

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
