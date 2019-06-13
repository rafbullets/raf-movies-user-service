package user

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Init - initialize router
func Init(r *mux.Router, db *sqlx.DB) {

	r.Methods("GET").Path("/api/users").HandlerFunc(Get(db))
	r.Methods("GET").Path("/api/users/{id}").HandlerFunc(GetOne(db))

	r.Methods("POST").Path("/api/users").HandlerFunc(Insert(db))
	r.Methods("POST").Path("/api/users/login").HandlerFunc(Login(db))
	r.Methods("POST").Path("/api/users/refresh").HandlerFunc(Refresh(db))

	r.Methods("PUT").Path("/api/users").HandlerFunc(Update(db))
}
