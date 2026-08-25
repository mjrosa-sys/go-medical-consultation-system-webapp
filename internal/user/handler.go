package user

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/render"
	"github.com/mjrosa-sys/go-medical-consultation-system-webapp/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type registerForm struct {
	Name     string
	Email    string
	Role     string
	Password string
	validator.Validator
}

type loginForm struct {
	Email    string
	Password string
	validator.Validator
}

type Handler struct {
	DB             *sql.DB
	SessionManager *scs.SessionManager
}

func NewHandler(db *sql.DB, sm *scs.SessionManager) *Handler {
	return &Handler{
		DB:             db,
		SessionManager: sm,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
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
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password cannot be more than 100 characters long")
	form.CheckField(validator.ValidPassword(form.Password), "password", "Password should contain at least one digit and character, and be 8 characters long")

	form.CheckField(validator.NotBlank(form.Role), "role", "A role should be defined")
	form.CheckField(validator.PermittedValue(form.Role, "patient", "doctor"), "role", "Invalid role")

	form.SetFormName("register")

	data := render.TemplateData{
		Validator: form.Validator,
		User: User{
			Name:  form.Name,
			Email: form.Email,
			Role:  form.Role,
		},
	}

	if !form.Valid() {
		ts, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/pages/home.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ts.ExecuteTemplate(w, "base", data)
		return
	}

	userModel := UserModel{DB: h.DB}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 14)
	if err != nil {
		log.Printf("Password hashing failed\n")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := userModel.Insert(form.Name, form.Email, form.Role, string(hashedPassword))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Printf("User created. ID: %d\n", id)

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	form := loginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.SetFormName("login")

	form.CheckField(validator.NotBlank(form.Email), "email", "Email field cannot be empty")
	form.CheckField(validator.MaxChars(form.Email, 256), "email", "Email cannot be more than 256 characters long")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Email format is not valid")

	form.CheckField(validator.NotBlank(form.Password), "password", "Password field cannot be empty")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password cannot be more than 100 characters long")
	form.CheckField(validator.ValidPassword(form.Password), "password", "Password should contain at least one digit and character, and be 8 characters long")

	// GetUserByEmail. Se não existe, adiciona o erro ao validator.Errors
	// bcrypt.Compare. Se retornar erro, adiciona ao validator.Errors

	userModel := UserModel{DB: h.DB}

	user, err := userModel.GetUserByEmail(form.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, ok := form.Errors["email"]; user == nil && !ok {
		form.AddFieldError("email", "Email not registered")
	}

	if user != nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
		if _, ok := form.Errors["password"]; err != nil && !ok {
			form.AddFieldError("password", "Wrong password")
		}
	}

	log.Printf("USER LOGIN: %v\n", form)

	data := render.TemplateData{
		Validator: form.Validator,
		User: User{
			Email: form.Email,
		},
	}

	if !form.Valid() {
		ts, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/pages/home.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserID", user.ID)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.SessionManager.Destroy(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
