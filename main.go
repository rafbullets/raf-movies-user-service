package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/user_service/api/user"
)

func main() {
	r := mux.NewRouter()

	db, err := sqlx.Connect("mysql", "Ney9CoabeM:ZFmcWctZXo@tcp(remotemysql.com)/Ney9CoabeM")
	if err != nil {
		panic(err)
	}

	user.Init(r, db)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, r))
}
