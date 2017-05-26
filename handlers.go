package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
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
		Role       string `json:"role"`
		IsDisabled bool   `json:"isdisabled"`
	}

	err := decodeJsonBody(w, r, &data)
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

	if data.Role != "admin" {
		data.Role = "user"
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

	u := User{
		Email:        data.Email,
		Username:     data.Username,
		Fullname:     data.Fullname,
		PasswordHash: string(hash),
		PasswordSalt: salt,
		Role:         data.Role,
		IsDisabled:   data.IsDisabled,
	}
	if err := db.Create(&u).Error; err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["status"] = "ok"

	resp.jSend(w)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := make(Response)

	if vars["email"] == "" {
		resp.jSendError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var u User
	fields := []string{"email", "username", "fullname", "role", "is_disabled"}
	if err := db.Select(fields).Where("email = ?", vars["email"]).First(&u).Error; err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["userinfo"] = u
	resp.jSend(w)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)

	u := make([]User, 10)
	fields := []string{"email", "username", "fullname", "role", "is_disabled"}
	if err := db.Select(fields).Find(&u).Error; err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["users"] = u
	resp.jSend(w)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	//Verify user is not already logged in
	resp := make(Response)

	session, err := store.Get(r, cookieName)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["logged"] == "true" {
		err = sessions.Save(r, w)
		if err != nil {
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp["status"] = "already logged"
		resp.jSend(w)
		return
	}
	//TODO define if old sessions should be updated in case session format has changed between logs

	//Read request parameters
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = decodeJsonBody(w, r, &data)
	if err != nil {
		return
	}

	//Validate data
	if data.Email == "" ||
		data.Password == "" {
		resp.jSendError(w, "Invalid attribute values", http.StatusBadRequest)
		return
	}

	var u User
	fields := []string{"username", "password_hash", "password_salt", "role", "is_disabled"}
	if err := db.Select(fields).Where("email = ?", data.Email).First(&u).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			resp.jSendError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		default:
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//TODO route to special page if user is disabled

	saltedPassword := u.PasswordSalt + data.Password

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(saltedPassword))
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
	session, err = store.Get(r, cookieName)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["accessLevel"] = u.Role
	session.Values["logged"] = "true"
	session.Values["language"] = "en-us"

	err = sessions.Save(r, w)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["status"] = "ok"
	resp.jSend(w)
}

func logOut(w http.ResponseWriter, r *http.Request) {
	resp := make(Response)
	session, err := store.Get(r, cookieName)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1

	err = sessions.Save(r, w)
	if err != nil {
		resp.jSendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["status"] = "ok"
	resp.jSend(w)
}
