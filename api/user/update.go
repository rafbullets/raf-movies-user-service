package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Update - update
func Update(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		decoder := json.NewDecoder(r.Body)

		var user User
		err := decoder.Decode(&user)
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
