package main

import (
	"crypto/rand"
	"encoding/base64"
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

	var data struct {
		Email      string `json:"email"`
		Username   string `json:"username"`
		Fullname   string `json:"fullname"`
		Password   string `json:"password"`
		IsDisabled bool   `json:"isdisabled"`
	}

	err := bodyToJson(w, r, &data)
	if err != nil {
		return
	}

	resp := make(Response)

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
		log.Printf("err: postUser: %s", err.Error())
		resp.jSendError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	saltedPassword := salt + data.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("err: postUser: %s", err.Error())
		resp.jSendError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

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
	vars := mux.Vars(r)
	resp := make(Response)

	session, err := store.Get(r, "jdata")
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["logged"] != "true" &&
		session.Values["accessLevel"] != "user" {
		resp.jSendError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

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
	err = db.QueryRow("select email, username, fullname, isdisabled from user_info where email = $1;",
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
}

func logIn(w http.ResponseWriter, r *http.Request) {
	//Read request parameters
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := bodyToJson(w, r, &data)
	if err != nil {
		return
	}

	resp := make(Response)

	//Validate data
	if data.Email == "" ||
		data.Password == "" {
		resp.jSendError(w, "Invalid attribute values", http.StatusBadRequest)
		return
	}

	//Process request
	var qResp struct {
		Username   string
		Hash       string
		Salt       string
		Isdisabled bool
	}
	err = db.QueryRow("select username, passwordhash, passwordsalt, isdisabled from user_info where email = $1;",
		data.Email,
	).Scan(&qResp.Username,
		&qResp.Hash,
		&qResp.Salt,
		&qResp.Isdisabled,
	)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//TODO route to special page if user is disabled

	saltedPassword := qResp.Salt + data.Password

	err = bcrypt.CompareHashAndPassword([]byte(qResp.Hash), []byte(saltedPassword))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.jSendError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		default:
			log.Printf("err: logIn: %s", err.Error())
			resp.jSendError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	session, err := store.Get(r, "jdata")
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["accessLevel"] = "user"
	session.Values["logged"] = "true"
	session.Values["language"] = "en-us"

	log.Printf("Session: %v", session)

	err = sessions.Save(r, w)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["status"] = "ok"
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
