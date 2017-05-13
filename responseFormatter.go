package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response map[string]interface{}

func (resp *Response) jSend(w http.ResponseWriter) {
	responseJson, err := json.Marshal(resp)
	if err != nil {
		responseCode := http.StatusInternalServerError
		responseJson = []byte(fmt.Sprintf("{ \"error\": \"%v\" }", err))
		w.WriteHeader(responseCode)
	}
	fmt.Fprintf(w, "%s", string(responseJson))
}

func (resp *Response) jSendError(w http.ResponseWriter, rErr string, responseCode int) {
	(*resp)["error"] = rErr
	w.WriteHeader(responseCode)
	resp.jSend(w)
}
