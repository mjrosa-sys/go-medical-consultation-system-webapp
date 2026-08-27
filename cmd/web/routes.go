package main

import (
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/appointment"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	appointmentHandler := appointment.NewHandler(app.db, app.sessionManager)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", app.redirectIfAuthenticated(http.HandlerFunc(home)))

	mux.HandleFunc("POST /register", app.UserRegister)
	mux.HandleFunc("POST /login", app.UserLogin)
	mux.HandleFunc("POST /logout", app.UserLogout)

	mux.Handle("GET /dashboard", app.requireAuthentication(http.HandlerFunc(dashboard)))

	mux.Handle("POST /appointments/create", app.requireAuthentication(http.HandlerFunc(appointmentHandler.Create)))

	return commonHeaders(app.sessionManager.LoadAndSave(mux))
}
