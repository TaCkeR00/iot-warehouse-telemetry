package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("missing token, query parameter"))
			return
		}

		_, err := app.auth.ValidateJWT(tokenString)
		if err != nil {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("invalid token %w", err))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) ValidateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("missing token, query parameter"))
			return
		}

		token, err := app.auth.ValidateJWT(tokenString)
		if err != nil {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("invalid token %w", err))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		isAdmin, ok := claims["is-admin"].(bool)
		if !ok || !isAdmin {
			app.httpError(w, http.StatusUnauthorized, fmt.Errorf("not admin"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
