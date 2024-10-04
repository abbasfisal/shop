package errors

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"shop/internal/pkg/sessions"
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

// SetErrors will return validation errors which translated into persian lang
func SetErrors(c *gin.Context, i18nBundle *i18n.Bundle, err error) {
	localizer := i18n.NewLocalizer(i18nBundle, "fa")

	switch err.(type) {
	case validator.ValidationErrors:
		errors := err.(validator.ValidationErrors)

		var errorMessages = make(map[string]string)
		for _, e := range errors {
			translation, _ := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: e.Tag(),
				TemplateData: map[string]interface{}{
					"Field": Mapper(strings.ToLower(e.Field())),
					"Param": e.Param(),
				},
			})
			errorMessages[strings.ToLower(e.Field())] = translation
		}
		marshal, _ := json.Marshal(errorMessages)

		sessions.Set(c, "errors", string(marshal))
		break

	default:
		sessions.Set(c, "errors", string(err.Error()))
	}

}

func Mapper(field string) string {
	switch field {
	case "firstname":
		return "نام"
	case "lastname":
		return "نام خانوادگی"
	case "mobile":
		return "موبایل"
	case "password":
		return "کلمه عبور"
	case "receiveraddress":
		return "آدرس تحویل گیرنده"
	case "receivermobile":
		return "موبایل"
	case "receivername":
		return "نام تحویل گیرنده"
	case "receiverpostalcode":
		return "کد پستی"

	default:
		return ""
	}
}
