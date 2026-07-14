package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	log.Println("Server running localhost:3001")

	err := http.ListenAndServe(":3001", mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}
