package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
)

//TODO load store keys from config file or env
var fs = "/Users/jrobin/Documents/jProjects/go/src/bitbucket.com/jrlangford/sessionsExample"
var store = sessions.NewFilesystemStore(fs, []byte("my-cookie-secret"))
var db *sql.DB
var r *mux.Router

func initRouter() {
	r = mux.NewRouter()

	r.HandleFunc("/cookie/save", saveSession)
	r.HandleFunc("/cookie/read", readSession)
	r.HandleFunc("/user", postUser).Methods("POST")
}

func safePing(db *sql.DB) {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("DB Ping Err: " + err.Error())
	}
}

func initDB() {
	log.Println("Connecting to database")
	var err error
	db, err = sql.Open("postgres", "host=localhost user=postgres dbname=postgres password=postgrespass sslmode=disable")
	if err != nil {
		log.Fatal("DB ERR: " + err.Error())
	}

	log.Println("Testing db connection")
	safePing(db)
	log.Println("Db connection sucessful")
}

func initSessionStore() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
}

func main() {
	initSessionStore()
	initDB()
	initRouter()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		err := db.Close()
		if err != nil {
			log.Print(err.Error())
			os.Exit(1)
		}
		log.Println("Db connection closed")
		os.Exit(0)
	}()

	http.ListenAndServe(":8080", r)
}
