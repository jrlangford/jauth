package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserInfo struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	Password   string `json:"password"`
	IsDisabled bool   `json:"isdisabled"`
}

const saltLength = 32

func generateSecureRandomBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(b)
	return out, nil
}

func createUser(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)

	decoder := json.NewDecoder(r.Body)
	var data UserInfo
	err := decoder.Decode(&data)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if data.Email == "" ||
		data.Username == "" ||
		data.Fullname == "" ||
		data.Password == "" {

		resp.jSendError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//TODO validate password

	salt, err := generateSecureRandomBytes(saltLength)
	if err != nil {
		// TODO: Do not propagate critical error data to client
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Salt: %v", salt)

	saltedPassword := salt + data.Password

	log.Printf("SaltedPass: %s", string(saltedPassword))

	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)

	var lastInsertId int
	err = db.QueryRow("insert into user_info (email, username, fullname, passwordhash, passwordsalt, isdisabled) values ($1, $2, $3, $4, $5, $6) returning id;",
		data.Email,
		data.Username,
		data.Fullname,
		hash,
		salt,
		data.IsDisabled,
	).Scan(&lastInsertId)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["id"] = lastInsertId

	resp.jSend(w)
}

func saveSession(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)
	resp["title"] = "This is part of the body"

	session, err := store.Get(r, "session-A")
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["val1"] = "This is value number one"
	session.Values["val2"] = "This is value number two"

	err = sessions.Save(r, w)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.jSend(w)
}

func readSession(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)

	session, err := store.Get(r, "session-A")
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["cookieVal1"] = session.Values["val1"].(string)
	resp["cookieVal2"] = session.Values["val2"].(string)
	resp.jSend(w)
}
