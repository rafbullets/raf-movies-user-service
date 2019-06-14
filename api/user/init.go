package user

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Init - initialize router
func Init(r *mux.Router, db *sqlx.DB) {

	r.Methods("GET").Path("/api/users").HandlerFunc(Get(db))
	r.Methods("GET").Path("/api/users/{id}").HandlerFunc(GetOne(db))

	r.Methods("POST").Path("/api/users/inc/{id}").HandlerFunc(Inc(db))
	r.Methods("POST").Path("/api/users").HandlerFunc(Insert(db))
	r.Methods("POST").Path("/api/users/login").HandlerFunc(Login(db))
	r.Methods("POST").Path("/api/users/refresh").HandlerFunc(Refresh(db))
	r.Methods("POST").Path("/api/users/ban/{id}").HandlerFunc(Ban(db))
	r.Methods("POST").Path("/api/users/unban/{id}").HandlerFunc(UnBan(db))

	r.Methods("PUT").Path("/api/users").HandlerFunc(Update(db))

	r.Methods("DELETE").Path("/api/users/{id}").HandlerFunc(Delete(db))
}
