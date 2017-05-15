package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

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

func postUser(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)

	decoder := json.NewDecoder(r.Body)
	var data struct {
		Email      string `json:"email,omitempty"`
		Username   string `json:"username,omitempty"`
		Fullname   string `json:"fullname,omitempty"`
		Password   string `json:"password,omitempty"`
		IsDisabled bool   `json:"isdisabled,omitempty"`
	}
	err := decoder.Decode(&data)
	if err != nil {
		switch err.Error() {
		case "EOF":
			resp.jSendError(w, "No body in request", http.StatusBadRequest)
			return
		default:
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer r.Body.Close()

	if data.Email == "" ||
		data.Username == "" ||
		data.Fullname == "" ||
		data.Password == "" {

		resp.jSendError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//TODO validate password complies with security standards

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

	_, err = db.Exec("insert into user_info (email, username, fullname, passwordhash, passwordsalt, isdisabled) values ($1, $2, $3, $4, $5, $6);",
		data.Email,
		data.Username,
		data.Fullname,
		hash,
		salt,
		data.IsDisabled,
	)
	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"user_info_email_key\"":
			resp.jSendError(w, "Email already registered", http.StatusConflict)
			return
		default:
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp["status"] = "ok"

	resp.jSend(w)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request) {

	//Read request parameters
	//Validate request
	//Process request
	//Respond

	vars := mux.Vars(r)
	resp := make(Response)

	if vars["email"] == "" {
		resp.jSendError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var qResp struct {
		Email      string `json:"email"`
		Username   string `json:"username"`
		Fullname   string `json:"fullname"`
		Isdisabled bool   `json:"isdisabled"`
	}
	err := db.QueryRow("select email, username, fullname, isdisabled from user_info where email = $1;",
		vars["email"],
	).Scan(&qResp.Email,
		&qResp.Username,
		&qResp.Fullname,
		&qResp.Isdisabled,
	)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["userinfo"] = qResp
	resp.jSend(w)

	//Query DB for user with given email

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
