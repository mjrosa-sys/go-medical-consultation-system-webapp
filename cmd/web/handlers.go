package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/models"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/render"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	ts.ExecuteTemplate(w, "base", nil)
}

func (app *application) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	userRole, err := app.userModel.GetUserRole(userID)
	if err != nil {
		http.Error(w, "Error when fetching user role", http.StatusInternalServerError)
	}

	data := render.TemplateData{
		User: models.User{
			Role: userRole,
		},
	}

	switch userRole {
	case "doctor":
		patients, err := app.userModel.GetAllPatients()
		if err != nil {
			log.Println(err)
			http.Error(w, "Error when fetching patients", http.StatusInternalServerError)
			return
		}

		aptmts, err := app.aptmtModel.GetAllAppointments()
		if err != nil {
			log.Println(err)
			http.Error(w, "Error when fetching appointments", http.StatusInternalServerError)
			return
		}

		data.Users = patients
		data.Appointments = aptmts

	case "patient":
		aptmts, err := app.aptmtModel.GetAppointmentsByUserId(userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error when fetching appointments", http.StatusInternalServerError)
			return
		}

		data.Appointments = aptmts

	default:
		log.Println("Unknown user role: ", err)
		http.Error(w, "Unknown user role", http.StatusInternalServerError)
		return
	}

	templates := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/dashboard.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	ts.ExecuteTemplate(w, "base", data)
}
