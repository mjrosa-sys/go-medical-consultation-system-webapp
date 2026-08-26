package appointment

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type NewAppointmentForm struct {
	PatientName        string
	AssignedDoctorName string
	DateAndTime        string
	Notes              string
	validator.Validator
}

type Handler struct {
	DB             *sql.DB
	SessionManager *scs.SessionManager
}

func NewHandler(DB *sql.DB, sm *scs.SessionManager) *Handler {
	return &Handler{
		DB:             DB,
		SessionManager: sm,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Form data: ", form)

	w.Write([]byte("Inserting a new appointment"))
}
