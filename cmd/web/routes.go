package main

import (
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/user"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := user.NewHandler(nil)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("POST /register", userHandler.Register)

	return mux
}
