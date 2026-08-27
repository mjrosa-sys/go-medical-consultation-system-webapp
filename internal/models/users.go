package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Role           string
	CreatedAt      time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, role, password string) (int, error) {
	stmt := "INSERT INTO users (name, email, hashed_password, role) VALUES (?, ?, ?, ?);"

	result, err := u.DB.Exec(stmt, name, email, password, role)
	if err != nil {
		log.Println("1", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("2", err)
		return 0, err
	}

	return int(id), nil
}

func (u *UserModel) GetUserByEmail(email string) (*User, error) {
	stmt := "select id, name, email, hashed_password, role, created_at from users where email = ?;"
	row := u.DB.QueryRow(stmt, email)

	user := User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) GetUserRole(ID int) (string, error) {
	stmt := `SELECT role
				FROM users
				WHERE id = ?;`

	row := u.DB.QueryRow(stmt, ID)

	var role string

	err := row.Scan(&role)
	if err != nil {
		return err.Error(), err
	}

	return role, nil
}
