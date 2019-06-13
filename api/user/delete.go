package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Delete - deletes user
func Delete(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		query := "Delete from user where id = ?"
		_, err := db.Exec(query, mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}

}
