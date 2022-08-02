package api

import (
	"log"
	"net/http"
)

func (s *server) middlewareJson(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

func (s *server) middlewareAccessLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next(w, r)
	}
}
