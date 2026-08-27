package main

import (
	"log"
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type NewAppointmentForm struct {
	PatientName        string
	AssignedDoctorName string
	DateAndTime        string
	Notes              string
	validator.Validator
}

func (app *application) AppointmentCreate(w http.ResponseWriter, r *http.Request) {
	//	doctorID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	//	if doctorID == 0 {
	//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//		return
	//	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed passing form", http.StatusInternalServerError)
	}

	form := NewAppointmentForm{
		PatientName:        r.PostForm.Get("patient_name"),
		AssignedDoctorName: r.PostForm.Get("doctor_name"),
		DateAndTime:        r.PostForm.Get("appointment_time"),
		Notes:              r.PostForm.Get("notes"),
	}

	form.SetFormName("NewAppointment")

	form.CheckField(validator.NotBlank(form.PatientName), "PatientName", "Patient name is required")
	form.CheckField(validator.MaxChars(form.PatientName, 100), "PatientName", "Patient name cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.AssignedDoctorName), "AssignedDoctorName", "Doctor name is required")
	form.CheckField(validator.MaxChars(form.AssignedDoctorName, 100), "AssignedDoctorName", "Doctor name cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.DateAndTime), "DateAndTime", "Date & Time is required")

	form.CheckField(validator.NotBlank(form.Notes), "Notes", "Notes is required")
	form.CheckField(validator.MaxChars(form.Notes, 256), "Notes", "Notes field cannot be more than 256 characters long")

	log.Println("Form data: ", form)
	log.Println("Form errors: ", form.Errors)

	w.Write([]byte("Inserting a new appointment"))
}
