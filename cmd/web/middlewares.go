package main

import (
	"log"
	"net/http"
)

func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authorizeDoctorUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doctorID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if doctorID == 0 {
			http.Error(w, "Forbidden access", http.StatusForbidden)
			return
		}

		userRole, err := app.userModel.GetUserRole(doctorID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error when fetching user role", http.StatusInternalServerError)
			return
		} else if userRole != "doctor" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.sessionManager.Exists(r.Context(), "authenticatedUserID") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate")

		next.ServeHTTP(w, r)
	})
}

func (app *application) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.sessionManager.Exists(r.Context(), "authenticatedUserID") {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
