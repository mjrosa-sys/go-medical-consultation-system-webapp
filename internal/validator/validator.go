package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FieldErrors map[string]string

type Validator struct {
	FormName string
	Errors   FieldErrors
}

func (v *Validator) SetFormName(name string) {
	v.FormName = name
}

func (v *Validator) GetFormName() string {
	return v.FormName
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddFieldError(key, errMessage string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = errMessage
	}
}

func (v *Validator) CheckField(ok bool, key, errMessage string) {
	if !ok {
		v.AddFieldError(key, errMessage)
	}
}

func NotBlank(fieldValue string) bool {
	return strings.TrimSpace(fieldValue) != ""
}

func MaxChars(fieldValue string, n int) bool {
	return utf8.RuneCountInString(fieldValue) <= n
}

func PermittedValue[T comparable](fieldValue T, PermittedValues ...T) bool {
	return slices.Contains(PermittedValues, fieldValue)
}

func ValidEmail(fieldValue string) bool {
	emailRX := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if len(fieldValue) > 254 {
		return false
	}
	return emailRX.MatchString(fieldValue)

}

func ValidPassword(fieldValue string) bool {
	var hasLetter, hasDigit, hasLength bool

	hasLength = len(fieldValue) >= 8

	for _, char := range fieldValue {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit && hasLength
}
