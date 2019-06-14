package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Insert - insert user
func Insert(db *sqlx.DB) http.HandlerFunc {
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

		decoder := json.NewDecoder(r.Body)

		var user User
		err = decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		query := "insert into user (last_name, first_name, email, username, role, number_of_reservations, password) values (?, ?, ?, ?, ?, ?, ?)"
		res, err := db.Exec(query, user.LastName, user.FirstName, user.Email, user.Username, user.Role, user.NumberOfReservations, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		id, _ := res.LastInsertId()

		user.ID = id

		json.NewEncoder(w).Encode(user)
	}
}

func Login(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")
		decoder := json.NewDecoder(r.Body)

		var user User
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.Get(&user, "SELECT * FROM user WHERE user.username = ? and user.password = ?", user.Username, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(3600 * time.Minute)

		var role string
		if user.Role == 0 {
			role = "admin"
		} else {
			role = "user"
		}
		claims := &Claims{
			ID:   user.ID,
			Role: role,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", "Bearer "+tokenString)
		json.NewEncoder(w).Encode(user)
	}
}

func Refresh(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
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

		// Now, create a new token for the current use, with a renewed expiration time
		expirationTime := time.Now().Add(3612 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", "Bearer "+tokenString)
	}
}

func Inc(db *sqlx.DB) http.HandlerFunc {
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

		id := mux.Vars(r)["id"]

		query := "Update user SET user.number_of_reservations = user.number_of_reservations + 1 where id = ?"
		_, err = db.Exec(query, id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func Ban(db *sqlx.DB) http.HandlerFunc {
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

		query := "Update user SET is_banned = ? where id = ?"
		_, err = db.Exec(query, true, mux.Vars(r)["id"])
		if err != nil {
			log.Fatal(err)
		}

		query = "insert into ban_activity (banned_user_id, banned_by_id) values (?, ?)"
		_, err = db.Exec(query, mux.Vars(r)["id"], claims.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

func UnBan(db *sqlx.DB) http.HandlerFunc {
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

		query := "Update user SET is_banned = ? where id = ?"
		_, err = db.Exec(query, false, mux.Vars(r)["id"])
		if err != nil {
			log.Fatal(err)
		}

		query = "Delete from ban_activity where banned_user_id = ?"
		_, err = db.Exec(query, mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}
