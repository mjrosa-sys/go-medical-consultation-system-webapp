package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/render"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
)

type registerForm struct {
	Name     string
	Email    string
	Role     string
	Password string
	validator.Validator
}

func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error passing form", http.StatusBadRequest)
		return
	}

	form := registerForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Role:     r.PostForm.Get("role"),
		Password: r.PostForm.Get("password"),
	}

	log.Println(form)

	form.CheckField(validator.NotBlank(form.Name), "name", "Name field cannot be empty")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "Name cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.Email), "email", "Email field cannot be empty")
	form.CheckField(validator.MaxChars(form.Email, 256), "email", "Email cannot be more than 256 characters long")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Email format is not valid")

	form.CheckField(validator.NotBlank(form.Password), "password", "Password field cannot be empty")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password cannot be more than 256 characters long")
	form.CheckField(validator.ValidPassword(form.Password), "password", "Password should contain at least one digit and character, and be 8 characters long")

	form.CheckField(validator.NotBlank(form.Role), "role", "A role should be defined")
	form.CheckField(validator.PermittedValue(form.Role, "patient", "doctor"), "role", "Invalid role")

	form.SetFormName("register")

	if !form.Valid() {
		data := render.TemplateData{
			Validator: form.Validator,
			User: User{
				Name:     form.Name,
				Email:    form.Email,
				Role:     form.Role,
				Password: form.Password,
			},
		}

		fmt.Println("TESTE: ", form.Email)

		ts, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/pages/home.tmpl")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ts.ExecuteTemplate(w, "base", data)
		return
	}

	w.Write([]byte("Creating new user"))
}
