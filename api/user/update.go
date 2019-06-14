package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

// Update - update
func Update(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Header.Get("Authorization")
		if c == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c = strings.Fields(c)[1]

		tknStr := c
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "application/json")
		decoder := json.NewDecoder(r.Body)

		var user User
		err = decoder.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		query := "Update user SET last_name = ?, first_name = ?, email = ?, username= ?, role = ?, number_of_reservations = ?, password = ? where id = ?"
		_, err = db.Exec(query, user.LastName, user.FirstName, user.Email, user.Username, user.Role, user.NumberOfReservations, user.Password, user.ID)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(user)

	}

}
