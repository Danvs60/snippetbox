package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// recover from panic with defer
		defer func(){
			// check if there has been a panic
			if err := recover(); err != nil {
				// this header tells Go mux to close the connection
				// or in HTTP/2 -> GOAWAY frame
				w.Header().Set("Connection", "close")
				// call server Error to feedback to user
				app.serverError(w, fmt.Errorf("%s:", err))
				// NOTE: why use fmt.Errorf -> recover returns an 'any' type
				// so we need to 'normalise' it by formatting to an error type
			}
		}()

		next.ServeHTTP(w, r)
	})
}
