package validator

import (
	"regexp"
	"slices"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) addError(key, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v Validator) Check(ok bool, key, value string) {
	if !ok {
		v.addError(key, value)
	}
}

func PermittedValue[T comparable](value T, PermittedValues ...T) bool {
	return slices.Contains(PermittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	unique := make(map[T]bool)

	for _, val := range values {
		unique[val] = true
	}
	return len(values) == len(unique)
}
