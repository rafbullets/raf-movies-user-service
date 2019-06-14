package status

import "github.com/dgrijalva/jwt-go"

type Status struct {
	ID         int64   `json:"id" db:"id"`
	FromPoints int64   `json:"fromPoints" db:"from_points"`
	ToPoints   int64   `json:"toPoints" db:"to_points"`
	Name       string  `json:"name" db:"name"`
	Discount   float64 `json:"discount" db:"discount"`
}

var JwtKey = []byte("my_secret_key")

type Claims struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}
