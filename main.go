package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/gorilla/mux"
	"upper.io/db.v3/mysql"
)

var config Config
var dbsession sqlbuilder.Database
var hashKey, blockKey []byte

func main() {
	var err error

	file, _ := os.Open("config.json")
	json.NewDecoder(file).Decode(&config)
	hashKey = []byte(config.HashKey)
	blockKey = []byte(config.BlockKey)

	var settings = mysql.ConnectionURL{
		Host:     config.MySQLHost,     // MySQL server IP or name.
		Database: config.MySQLDatabase, // Database name.
		User:     config.MySQLUser,     // Optional user name.
		Password: config.MySQLPassword, // Optional user password.
	}

	dbsession, err = mysql.Open(settings)

	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}

	defer dbsession.Close()

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/board/create", CreateBoard).Methods("POST")
	api.HandleFunc("/board/list", ListBoard).Methods("GET")
	api.HandleFunc("/board/{id}", DeleteBoard).Methods("DELETE")
	api.HandleFunc("/board/{id}", UpdateBoard).Methods("POST")
	api.HandleFunc("/board/{id}/{start}/{amount}", ListThread).Methods("GET")

	api.HandleFunc("/thread/{id}/create", CreateThread).Methods("POST")
	api.HandleFunc("/thread/list/{id}/{start}/{amount}", ListThread).Methods("GET")
	api.HandleFunc("/thread/{id}/{start}/{amount}", GetThread).Methods("GET")
	api.HandleFunc("/thread/{id}", DeleteThread).Methods("DELETE")
	api.HandleFunc("/thread/{id}", UpdateThread).Methods("POST")

	api.HandleFunc("/post/{id}/create", CreatePost).Methods("POST")
	api.HandleFunc("/post/list/{id}/{start}/{amount}", GetThread).Methods("GET")
	api.HandleFunc("/post/{id}", GetPost).Methods("GET")
	api.HandleFunc("/post/{id}", DeletePost).Methods("DELETE")
	api.HandleFunc("/post/{id}", UpdatePost).Methods("POST")

	api.HandleFunc("/user/code", GetCode).Methods("GET")
	api.HandleFunc("/user/code", SetCode).Methods("POST")
	api.HandleFunc("/user/code", DeleteCode).Methods("DELETE")

	api.HandleFunc("/setup", SetupBBS).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8443",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
