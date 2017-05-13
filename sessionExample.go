package main

import (
	"database/sql"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//TODO load store keys from config file or env
var fs = "/Users/jrobin/Documents/jProjects/go/src/bitbucket.com/jrlangford/sessionsExample"
var store = sessions.NewFilesystemStore(fs, []byte("my-cookie-secret"))
var db *sql.DB

func setRoutes() {
	http.HandleFunc("/cookie/save", saveSession)
	http.HandleFunc("/cookie/read", readSession)
	http.HandleFunc("/user/create", createUser)
}

func safePing(db *sql.DB) {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("DB Ping Err: " + err.Error())
	}
}

func main() {
	//Set oprions for session store
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	log.Println("Connecting to database")
	var err error
	db, err = sql.Open("postgres", "host=localhost user=postgres dbname=postgres password=postgrespass sslmode=disable")
	if err != nil {
		log.Fatal("DB ERR: " + err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}()

	log.Println("Testing db connection")
	safePing(db)
	log.Println("Db connection sucessful")

	setRoutes()
	//Wrap handlers with context.ClearHandler to avoid leaking memory with gorilla/sessions
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
