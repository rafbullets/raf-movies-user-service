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

		var status Status
		err = decoder.Decode(&status)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(status)
		query := "Update status SET from_points = ?, to_points = ?, name = ?, discount= ? where id = ?"
		_, err = db.Exec(query, status.FromPoints, status.ToPoints, status.Name, status.Discount, status.ID)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(status)

	}

}
