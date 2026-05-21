package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		// parse to get token
		// Authorization: Basic <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorized header is malformed"))
			return
		}

		// decode token
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedBasicErrorResponse(w, r, err)
			return
		}

		username := app.config.auth.basic.user
		password := app.config.auth.basic.pass

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
			return
		}
		// ctx := context.WithValue(r.Context(), "user", "123")

		app.logger.Info("successfully logged in basic auth")
		next.ServeHTTP(w, r)
	})
}
