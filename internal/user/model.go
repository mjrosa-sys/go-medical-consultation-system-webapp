package user

import (
	"database/sql"
	"log"
)

type User struct {
	Name     string
	Email    string
	Role     string
	Password string
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
