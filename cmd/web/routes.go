package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", app.redirectIfAuthenticated(http.HandlerFunc(app.Home)))

	mux.HandleFunc("POST /users/register", app.UserRegister)
	mux.HandleFunc("POST /users/login", app.UserLogin)
	mux.HandleFunc("POST /users/logout", app.UserLogout)

	mux.Handle("GET /dashboard", app.requireAuthentication(http.HandlerFunc(app.Dashboard)))

	mux.Handle("POST /appointments/create", app.requireAuthentication(http.HandlerFunc(app.AppointmentCreate)))

	return commonHeaders(app.sessionManager.LoadAndSave(mux))
}
