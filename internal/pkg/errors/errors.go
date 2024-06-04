package errors

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
)

var errorList = make(map[string]string)

func Init() {
	errorList = map[string]string{}
}

func SetFromErrors(err error) {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			Add(fieldErr.Field(), GetErrorMessage(fieldErr.Tag()))
		}
	}
}
func Add(key string, value string) {
	errorList[strings.ToLower(key)] = value
}

func GetErrorMessage(tag string) string {
	return errorMessages()[tag]
}

func Get() map[string]string {
	return errorList
}
func ToString() string {
	out, _ := json.Marshal(Get())

	if out != nil {
		return string(out)
	}
	return ""
}

func errorMessages() map[string]string {
	return map[string]string{
		"required": "The field is required",
		"email":    "The field must have valid email",
		"min":      "Should be more than limit",
		"max":      "Should be less than limit",
	}
}

type MyError struct {
	MainMessage string
	Message     string
	Code        int
}

func NewMyError(mainMessage string, message string, code int) MyError {
	return MyError{
		MainMessage: mainMessage,
		Message:     message,
		Code:        code,
	}
}
func (m MyError) Error() string {
	return m.Message
}
