package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/models"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/render"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type NewAppointmentForm struct {
	PatientID   int
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
		http.Error(w, "Failed parsing form", http.StatusBadRequest)
		return
	}

	// Convert form input "patient_id" (string) to integer
	patientID, err := strconv.Atoi(r.PostForm.Get("patient_id"))
	if err != nil {
		patientID = 0 // Will fail the validation check below
	}

	form := NewAppointmentForm{
		PatientID:   patientID,
		DateAndTime: r.PostForm.Get("appointment_time"),
		Notes:       r.PostForm.Get("notes"),
	}

	form.SetFormName("NewAppointment")

	parsedTime, err := time.Parse("2006-01-02T15:04", form.DateAndTime)
	if err != nil {
		form.AddFieldError("DateAndTime", "Invalid date format")
	}

	form.CheckField(validator.NotBlank(form.DateAndTime), "DateAndTime", "Date & Time is required")

	// Validate PatientID is selected
	form.CheckField(form.PatientID > 0, "PatientID", "Please select a valid patient")

	form.CheckField(validator.NotBlank(form.Notes), "Notes", "Notes is required")
	form.CheckField(validator.MaxChars(form.Notes, 256), "Notes", "Notes field cannot be more than 256 characters long")

	log.Println("Form data: ", form)
	log.Println("Form errors: ", form.Errors)

	if !form.Valid() {
		// Fetch patients again so the template can re-render the dropdown options on validation error
		patients, err := app.userModel.GetAllPatients() // Adjust model reference as needed
		if err != nil {
			log.Println("Failed to retrieve patients on validation error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := render.TemplateData{
			User:      models.User{Role: "doctor"},
			Users:     patients,
			Validator: form.Validator,
		}

		ts, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/pages/dashboard.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ts.ExecuteTemplate(w, "base", data)
		return
	}

	// Call updated Insert method with PatientID
	_, err = app.aptmtModel.Insert(form.PatientID, doctorID, form.Notes, parsedTime)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error when inserting new appointment in the database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
