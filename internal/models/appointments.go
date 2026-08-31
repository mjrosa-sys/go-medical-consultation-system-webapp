package models

import (
	"database/sql"
	"log"
	"time"
)

type Appointment struct {
	ID                 int
	PatientName        string
	DoctorID           int
	AssignedDoctorName string
	Notes              string
	DateAndTime        time.Time
	CreatedAt          time.Time
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

func (aptmt *AppointmentModel) GetAllAppointments() ([]Appointment, error) {
	stmt := `
		SELECT 
			a.id, 
			a.patient_name, 
			a.doctor_id, 
			u.name AS assigned_doctor_name, 
			COALESCE(a.notes, '') AS notes, 
			a.date_and_time, 
			a.created_at
		FROM appointments a
		INNER JOIN users u ON a.doctor_id = u.id
		ORDER BY a.date_and_time ASC;`

	rows, err := aptmt.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment

		err := rows.Scan(
			&a.ID,
			&a.PatientName,
			&a.DoctorID,
			&a.AssignedDoctorName,
			&a.Notes,
			&a.DateAndTime,
			&a.CreatedAt,
		)
		if err != nil {
			log.Println("Scan error on appointment row:", err)
			return nil, err
		}

		appointments = append(appointments, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}
