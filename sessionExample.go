package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Title string `json:"title,omitempty"`
}

var store = sessions.NewCookieStore([]byte("my-cookie-secret"))

func (resp *Response) send(w http.ResponseWriter) {
	responseJson, err := json.Marshal(resp)
	if err != nil {
		responseCode := http.StatusInternalServerError
		responseJson = []byte(fmt.Sprintf("{ \"error\": \"%v\" }", err))
		w.WriteHeader(responseCode)
	}
	fmt.Fprintf(w, "%s", string(responseJson))
}

func (resp *Response) sendError(w http.ResponseWriter, rErr string, responseCode int) {
	resp.Error = rErr
	w.WriteHeader(responseCode)
	resp.send(w)
}

func saveSession(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Title: "This is part of the body",
	}

	session, err := store.Get(r, "session-A")
	if err != nil {
		resp.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["val1"] = "This is value number one"
	session.Values["val2"] = "This is value number two"

	err = sessions.Save(r, w)
	if err != nil {
		resp.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.send(w)
}

func readSession(w http.ResponseWriter, r *http.Request) {
	resp := Response{}

	session, err := store.Get(r, "session-A")
	if err != nil {
		resp.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Title = session.Values["val1"].(string)
	resp.send(w)
}

func setRoutes() {
	http.HandleFunc("/cookie/save", saveSession)
	http.HandleFunc("/cookie/read", readSession)
}

func main() {
	setRoutes()
	//Wrap handlers with context.ClearHandler to avoid leaking memory
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
