// copywrite 2022 Matthew R Kasun mkasun@nusak.ca
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetIP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RemoteAddr))
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("login")
	var err error
	request := struct {
		User string
		Pass string
	}{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid request "+err.Error())
		return
	}
	if err := json.Unmarshal(body, &request); err != nil {
		fail(w, http.StatusBadRequest, "invalid request "+err.Error())
		return
	}
	if request.User != "demo" || request.Pass != "pass" {
		fail(w, http.StatusBadRequest, "invalid username or password")
		return
	}
	expires := time.Now().Add(time.Minute * 3)
	claims := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "nusak.ca",
			ID:        request.User,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JWT, err := token.SignedString(SigningKey)
	if err != nil {
		fail(w, http.StatusInternalServerError, "unable to create JWT "+err.Error())
		return
	}
	response := struct {
		JWT string
	}{
		JWT: JWT,
	}
	payload, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	GetIP(w, r)
}
