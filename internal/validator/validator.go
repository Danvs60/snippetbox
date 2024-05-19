package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator map for validation errors for forms
type Validator struct {
	FieldErrors map[string]string
}

// Form is valid if FieldErrors is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds error message to FieldErrors (if it does not already exist)
func (v *Validator) AddFieldError(key, message string) {
	// initialise map if needed
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}
