package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response map[string]interface{}

func (resp *Response) jSend(w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		responseCode := http.StatusInternalServerError
		w.WriteHeader(responseCode)
		log.Printf("err: jSend: %s", err.Error())
	}
}

func (resp *Response) jSendError(w http.ResponseWriter, rErr string, responseCode int) {
	(*resp)["error"] = rErr
	w.WriteHeader(responseCode)
	resp.jSend(w)
}

func decodeJsonBody(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err = decoder.Decode(&data)
	if err != nil {
		resp := make(Response)
		switch err.Error() {
		case "EOF":
			resp.jSendError(w, "No body in request", http.StatusBadRequest)
			return
		default:
			log.Printf("err: bodyToJson: %s", err.Error())
			resp.jSendError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	return
}
