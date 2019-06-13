package main

import (
	"fmt"
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
	log.Fatal(http.ListenAndServe(GetPort(), r))
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
