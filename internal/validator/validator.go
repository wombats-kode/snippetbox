package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile() function to sanity check an email address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Add a new NonFieldErrors []string field to hold errors that are not associated
// related to specific form fields
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() returns true if the Field Errors map doesnt contain any entries
// Check for NonFieldErrors also
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Add a new NonFieldErrors []string function to add error messages
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError() adds an error message to the field errors map (so long
// as no entriy already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note we need to initialise the map first if isnt already
	// initialised
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldError map only if a
// Validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue returns true if the type T equals one of the permitted variadic permittedValues.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// MinChars() returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches() returns true if the value matches a provided compiled regular
// expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
