package status

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

func Insert(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		fmt.Println("EVO ME")
		c := r.Header.Get("Authorization")
		if c == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println("EVO ME")

		c = strings.Fields(c)[1]
		fmt.Println("EVO ME")

		tknStr := c
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		fmt.Println("EVO ME")

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println("EVO ME")

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("EVO ME")

		decoder := json.NewDecoder(r.Body)

		var status Status
		err = decoder.Decode(&status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("EVO ME")

		query := "insert into status (from_points, to_points, name, discount) values (?, ?, ?, ?)"
		res, err := db.Exec(query, status.FromPoints, status.ToPoints, status.Name, status.Discount)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("EVO ME")

		id, _ := res.LastInsertId()

		status.ID = id

		json.NewEncoder(w).Encode(status)
		fmt.Println("EVO ME")

	}
}
