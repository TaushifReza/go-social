package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "authorization header is missing", "authorization header is missing")
				return
			}
			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				writeUnauthorized(w, "authorization header is malformed", "authorization header is malformed")
				return
			}
			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				writeUnauthorized(w, "authorization header is malformed", "authorization header is malformed")
				return
			}
			creds := strings.SplitN(string(decoded), ":", 2)
			// check credentials
			if len(creds) != 2 || creds[0] != "admin" || creds[1] != "admin" {
				writeUnauthorized(w, "invalid credentials", "invalid credentials")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
