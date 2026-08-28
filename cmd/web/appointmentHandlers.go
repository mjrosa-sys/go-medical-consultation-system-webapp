package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/models"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/render"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type NewAppointmentForm struct {
	PatientName string
	DateAndTime string
	Notes       string
	validator.Validator
}

func (app *application) AppointmentCreate(w http.ResponseWriter, r *http.Request) {
	doctorID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if doctorID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("User ID: ", doctorID)

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed passing form", http.StatusInternalServerError)
	}

	form := NewAppointmentForm{
		PatientName: r.PostForm.Get("patient_name"),
		DateAndTime: r.PostForm.Get("appointment_time"),
		Notes:       r.PostForm.Get("notes"),
	}

	form.SetFormName("NewAppointment")

	form.CheckField(validator.NotBlank(form.PatientName), "PatientName", "Patient name is required")
	form.CheckField(validator.MaxChars(form.PatientName, 100), "PatientName", "Patient name cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.DateAndTime), "DateAndTime", "Date & Time is required")

	form.CheckField(validator.NotBlank(form.Notes), "Notes", "Notes is required")
	form.CheckField(validator.MaxChars(form.Notes, 256), "Notes", "Notes field cannot be more than 256 characters long")

	log.Println("Form data: ", form)
	log.Println("Form errors: ", form.Errors)

	data := render.TemplateData{
		User:      models.User{Role: "doctor"},
		Validator: form.Validator,
	}

	if !form.Valid() {
		ts, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/pages/dashboard.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ts.ExecuteTemplate(w, "base", data)
		return
	}

	// Insere no DB
	// Redireciona para /dashboard

	w.Write([]byte("Inserting a new appointment"))
}
