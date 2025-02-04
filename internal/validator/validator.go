package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// FieldName represents that name of the field whith errors, like "title",
// "content" or "expires".
type FieldName = string

// FieldErrorMessage represents the error description associated with the
// given field name.
type FieldErrorMessage = string

type Validator struct {
	FieldErrors map[FieldName]FieldErrorMessage
}

// Valid returns true if there are no errors recorded in any of the fields;
// false otherwise.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds the error message to the given field if there is no
// existing message for that field.
func (v *Validator) AddFieldError(name FieldName, message FieldErrorMessage) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[FieldName]FieldErrorMessage)
	}

	if _, exists := v.FieldErrors[name]; !exists {
		v.FieldErrors[name] = message
	}
}

// CheckField adds an error message to the given field name if the validation
// check is not ok.
func (v *Validator) CheckField(ok bool, name FieldName, message FieldErrorMessage) {
	if !ok {
		v.AddFieldError(name, message)
	}
}

// NotBlank returns a true if the string input is not an empty string; false
// otherwise.
func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if the given string contains no more than max
// characters; false otherwise.
func (v *Validator) MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue returns true if the given value is one of the permitted
// values; false otherwise.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
