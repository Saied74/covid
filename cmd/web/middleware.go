//adding request logging and panic reco very

package main

import (
	"fmt"
	"net/http"
)

func (st *sT) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
			r.URL.RequestURI())

		next.ServeHTTP(w, r)

	})
}

func (st *sT) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//only runs if there is a panic
		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				st.serverError(w, fmt.Errorf("%s", err))

			}
		}()
		next.ServeHTTP(w, r)
	})
}
