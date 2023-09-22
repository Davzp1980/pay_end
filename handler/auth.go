package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"pay/internal"
	"pay/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("My_key")

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isadmin"`
	jwt.RegisteredClaims
}

func SignIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input internal.Input

		json.NewDecoder(r.Body).Decode(&input)

		userName, err := repository.GetUser(db, input.Name, input.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)

		claims := &Claims{
			Username: userName,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Welcome %s", userName)))

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
