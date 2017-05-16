package main

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/boj/redistore.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var store *redistore.RediStore
var db *sql.DB
var r *mux.Router

const cookieName = "jdata"

func initRouter() {

	usersAndAdmins := []string{"user", "admin"}
	//users := []string{"user"}
	admins := []string{"admin"}

	r = mux.NewRouter()

	logH := handlers.LoggingHandler(os.Stdout, r)
	uaH := auth(logH, usersAndAdmins)
	aH := auth(logH, admins)

	uaCorsH := handlers.CORS()(uaH)
	aCorsH := handlers.CORS()(aH)
	//TODO let each resource define who has acces to it instead of defining it on a route level

	http.Handle("/", logH)

	//Test
	r.HandleFunc("/cookie/save", saveSession)
	r.HandleFunc("/cookie/read", readSession)

	//User
	r.HandleFunc("/users", postUser).Methods("POST")

	//Access
	r.HandleFunc("/login", logIn).Methods("POST")

	http.Handle("/logout", uaCorsH)
	r.HandleFunc("/logout", logOut).Methods("POST")

	//Admin
	http.Handle("/admins/", aCorsH)
	admin := r.PathPrefix("/admins").Subrouter()
	admin.HandleFunc("/users/{email}", getUserByEmail).Methods("GET")
	admin.HandleFunc("/users", getUsers).Methods("GET")

}

func safePing(db *sql.DB) {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("DB Ping Err: " + err.Error())
	}
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", "host=localhost user=postgres dbname=postgres password=postgrespass sslmode=disable")
	if err != nil {
		log.Fatal("DB ERR: " + err.Error())
	}
	safePing(db)
	log.Println("Db connection sucessful")
}

func initSessionStore() {
	var err error
	store, err = redistore.NewRediStore(10, "tcp", ":8000", "", []byte("secret-key"))
	if err != nil {
		log.Fatal("Store ERR: " + err.Error())
	}
	log.Println("Store connection sucessful")

	const secondsInDay = 86400

	store.SetMaxAge(secondsInDay * 7)
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

		err = store.Close()
		if err != nil {
			log.Print(err.Error())
			os.Exit(1)
		}
		log.Println("Store connection closed")
		os.Exit(0)
	}()

	http.ListenAndServe(":8080", nil)
}
