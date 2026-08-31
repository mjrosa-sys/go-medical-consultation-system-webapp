package models

import (
	"database/sql"
	"log"
	"time"
)

type Appointment struct {
	ID                 int
	PatientID          int
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

func (aptmt *AppointmentModel) Insert(patientID int, doctorID int, notes string, dateAndTime time.Time) (int, error) {
	stmt := `INSERT INTO appointments (patient_id, doctor_id, notes, date_and_time)
            VALUES (?, ?, ?, ?);`

	result, err := aptmt.DB.Exec(stmt, patientID, doctorID, notes, dateAndTime)
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
			a.patient_id, 
			p.name AS patient_name, 
			a.doctor_id, 
			d.name AS assigned_doctor_name, 
			COALESCE(a.notes, '') AS notes, 
			a.date_and_time, 
			a.created_at
		FROM appointments a
		INNER JOIN users p ON a.patient_id = p.id
		INNER JOIN users d ON a.doctor_id = d.id
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
			&a.PatientID,
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

func (aptmt *AppointmentModel) GetAppointmentsByUserId(id int) ([]Appointment, error) {
	stmt := `
		SELECT 
			a.id, 
			a.patient_id, 
			p.name AS patient_name, 
			a.doctor_id, 
			d.name AS assigned_doctor_name, 
			COALESCE(a.notes, '') AS notes, 
			a.date_and_time, 
			a.created_at
		FROM appointments a
		INNER JOIN users p ON a.patient_id = p.id
		INNER JOIN users d ON a.doctor_id = d.id
		WHERE a.patient_id = ? OR a.doctor_id = ?
		ORDER BY a.date_and_time ASC;`

	rows, err := aptmt.DB.Query(stmt, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment

		err := rows.Scan(
			&a.ID,
			&a.PatientID,
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
