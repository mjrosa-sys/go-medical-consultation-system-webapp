package main

import (
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/user"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	userHandler := user.NewHandler(app.db, app.sessionManager)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)

	mux.HandleFunc("POST /register", userHandler.Register)

	mux.HandleFunc("GET /dashboard", dashboard)

	return commonHeaders(app.sessionManager.LoadAndSave(mux))
}
