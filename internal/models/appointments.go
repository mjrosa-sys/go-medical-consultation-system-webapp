package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	ID          int
	PatientName string
	DoctorID    int
	Notes       string
	DateAndTime time.Time
	CreatedAt   time.Time
}

type AppointmentModel struct {
	DB *sql.DB
}

func (aptmt *AppointmentModel) Insert(patientName string, doctorID int, notes string, dateAndTime time.Time) (int, error) {
	stmt := `INSERT INTO
				appointments (patient_name, doctor_id, notes, date_and_time)
				VALUES (?, ?, ?, ?);`

	result, err := aptmt.DB.Exec(stmt, patientName, doctorID, notes, dateAndTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
