package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/boj/redistore.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var store *redistore.RediStore
var db *gorm.DB
var r *mux.Router

const cookieName = "jdata"

func safePing(db *gorm.DB) {
	err := db.Exec("SELECT 1").Error
	if err != nil {
		log.Fatal("DB Ping Err: " + err.Error())
	}
}

func initGDB() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=postgres password=postgrespass sslmode=disable")
	if err != nil {
		log.Fatal("DB ERR: " + err.Error())
	}
	safePing(db)
	log.Println("Db connection sucessful")

	db.AutoMigrate(&User{})
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
	initGDB()
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
		log.Println("Gdb connection closed")

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
