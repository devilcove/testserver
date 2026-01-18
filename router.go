package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/devilcove/mux"
	"github.com/golang-jwt/jwt/v4"
)

var SigningKey []byte

type CustomClaims struct {
	jwt.RegisteredClaims
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter(mux.Logger)
	// basic routes
	router.Get("/ip", GetIP)
	router.Post("/login", Login)
	// authenticated routes
	restricted := router.Group("/api", Auth)
	//{
	restricted.Get("/hello", Hello)
	//}
	return router
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header["Authorization"]) == 0 {
			fail(w, http.StatusUnauthorized, "no auth header")
			return
		}
		id, status := getFromJWT(r.Header["Authorization"][0])
		log.Println(id, status, time.Now())
		if status == 1 {
			fail(w, http.StatusUnauthorized, "token expired")
			return
		}
		if status == 2 {
			fail(w, http.StatusUnauthorized, "invalid token")
			return
		}
		if id != "demo" {
			fail(w, http.StatusUnauthorized, "no such user")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func fail(w http.ResponseWriter, status int, message string) {
	log.Output(2, message)
	response := struct {
		Message string
	}{
		Message: message,
	}
	payload, _ := json.Marshal(response)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func getFromJWT(auth string) (string, int) {
	log.Println(auth)
	token, err := jwt.ParseWithClaims(auth, &CustomClaims{}, func(*jwt.Token) (any, error) {
		return SigningKey, nil
	})
	if err != nil {
		log.Println("error from jwt parse ", err)
		return "", 1
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Println("claims ", claims, "token ", token.Valid)
		return claims.ID, 0
	}
	return "", 2
}
