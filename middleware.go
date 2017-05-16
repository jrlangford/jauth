package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func auth(next http.Handler, accessLevels []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := make(Response)

		session, err := store.Get(r, "jdata")
		if err != nil {
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session.Values["logged"] != "true" {
			resp.jSendError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		levelMatch := false
		for _, al := range accessLevels {
			if session.Values["accessLevel"] == al {
				levelMatch = true
			}
		}
		if !levelMatch {
			resp.jSendError(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		//Renew session's ttl
		err = sessions.Save(r, w)
		if err != nil {
			resp.jSendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
