package status

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Init - initialize router
func Init(r *mux.Router, db *sqlx.DB) {

	r.Methods("GET").Path("/api/status").HandlerFunc(Get(db))
	r.Methods("POST").Path("/api/status").HandlerFunc(Insert(db))

	r.Methods("PUT").Path("/api/status").HandlerFunc(Update(db))

	r.Methods("DELETE").Path("/api/status/{id}").HandlerFunc(Delete(db))
}
