package main

import (
	"log"
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/user"
)

func main() {
	mux := http.NewServeMux()

	userHandler := user.NewHandler(nil)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("POST /register", userHandler.Register)

	log.Println("Server running localhost:3001")

	err := http.ListenAndServe(":3001", mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}
