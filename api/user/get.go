package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

// Get - gets all users from database
func Get(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		users := []User{}
		err := db.Select(&users, "SELECT * FROM user")
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

// Get - specific user
func GetOne(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		user := User{}
		err := db.Get(&user, "SELECT * FROM user WHERE user.id = ?", mux.Vars(r)["id"])
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(user)
	}
}
