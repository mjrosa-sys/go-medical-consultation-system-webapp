package validator

import (
	"net/mail"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, errMessage string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = errMessage
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
	_, err := mail.ParseAddress(fieldValue)
	return err == nil
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
