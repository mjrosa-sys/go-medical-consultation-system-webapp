package render

import (
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/models"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type TemplateData struct {
	Validator    validator.Validator
	User         models.User
	Users        []models.User
	Appointment  models.Appointment
	Appointments []models.Appointment
}
