package rest

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const allowedOrigin = "*"

//Auth middleware recupera token do header, valida e adiciona no header hotel_id
func Auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", allowedOrigin)
		userID := r.Header.Get("Authorization")
		r.Header.Add("user_id", userID)
		inner.ServeHTTP(w, r)
	})
}

//Cors middleware retorna pre-fligth OPTIONS
func Cors(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, Access-Control-Allow-Origin")
			w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, GET, POST")

			w.WriteHeader(http.StatusOK)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"LogMW %s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
