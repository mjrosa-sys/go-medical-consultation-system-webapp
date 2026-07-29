package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("Server running localhost:3001")

	err := http.ListenAndServe(":3001", routes())
	if err != nil {
		log.Fatal(err.Error())
	}
}
