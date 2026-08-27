package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	ID          int
	PatientName string
	DoctorID    int
	DateAndTime time.Time
	CreatedAt   time.Time
}

type AppointmentModel struct {
	DB *sql.DB
}

func (appointment *AppointmentModel) Insert(patientName string, doctorID int, dateAndTime, createdAt time.Time) (int, error) {

	return 0, nil
}
