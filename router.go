package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func initRouter() {

	r = mux.NewRouter()

	logH := handlers.LoggingHandler(os.Stdout, r)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8090"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	cors := handlers.CORS(originsOk, headersOk, methodsOk)

	public := cors(logH)
	usersAndAdmins := cors(auth([]string{"user", "admin"}, logH))
	admins := cors(auth([]string{"admin"}, logH))

	http.Handle("/", public)
	r.HandleFunc("/users", postUser).Methods("POST")
	r.HandleFunc("/login", logIn).Methods("POST")

	http.Handle("/logout", usersAndAdmins)
	r.HandleFunc("/logout", logOut).Methods("POST")

	http.Handle("/admins/", admins)
	admin := r.PathPrefix("/admins").Subrouter()
	admin.HandleFunc("/users/{email}", getUserByEmail).Methods("GET")
	admin.HandleFunc("/users", getUsers).Methods("GET")
}
