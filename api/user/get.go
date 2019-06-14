package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/user_service/api/status"
)

// Get - gets all users from database
func Get(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

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

		users := []User{}
		page, byPage, sort_by, sort_order, filter_field, filter_value := r.URL.Query()["page"][0], r.URL.Query()["byPage"][0], r.URL.Query()["sort_by"][0], r.URL.Query()["sort_order"][0],
			r.URL.Query()["filter_field"][0], r.URL.Query()["filter_value"][0]
		query := "SELECT * FROM user"

		if filter_value != "" && filter_field != "" {
			query += " WHERE " + filter_field + " LIKE '%" + filter_value + "%' "
		}
		if sort_by != "" {
			query += " ORDER BY " + sort_by
		}

		if sort_order != "" {
			query += " " + sort_order
		}
		if page != "" && byPage != "" {
			pageInt, _ := strconv.Atoi(page)
			byPageInt, _ := strconv.Atoi(byPage)
			offset := pageInt * byPageInt
			query += " LIMIT " + strconv.Itoa(byPageInt) + " OFFSET " + strconv.Itoa(offset)
		}

		err = db.Select(&users, query)
		if err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

// Get - specific user
func GetOne(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

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

		user := User{}
		err = db.Get(&user, "SELECT * FROM user WHERE user.id = ? ", mux.Vars(r)["id"])
		if err != nil {
			log.Println("EVO ME")
			log.Println(err)
			return
		}

		s, err := status.GetOneByPoints(db, user.NumberOfReservations)
		if err != nil {
			log.Println(err)
			return
		}

		result := UserStatus{user, *s}

		json.NewEncoder(w).Encode(result)
	}
}
