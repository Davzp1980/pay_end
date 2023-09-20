package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

func UserIdentification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sign-in" && r.URL.Path != "/create-admin" && r.URL.Path != "/sign-up" {

			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tokenString := c.Value

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		}
		next.ServeHTTP(w, r)
	})
}
