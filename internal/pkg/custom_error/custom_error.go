package custom_error

import (
	"errors"
	"gorm.io/gorm"
)

type CustomError struct {
	OriginalMessage string
	DisplayMessage  string //will set to user
	Code            int
}

func New(original, display string, code int) CustomError {
	return CustomError{
		OriginalMessage: original,
		DisplayMessage:  display,
		Code:            code,
	}
}
func (ce CustomError) Error() string {
	return ce.DisplayMessage
}

const (
	MobileOrPasswordIsWrong string = "شماره موبایل یا رمزعبور اشتباه است"
	InternalServerError     string = "خطای سرور"
	RecordNotFound          string = "رکوردی یافت نشد"
	MustBeUnique            string = "مقدار وارد شده باید منحصربه فرد باشد"
	IsRequired              string = "فیلد مورد نظر اجباری است"
	MustBeImage             string = "پسوند عکس باید jpg|jpeg|png باشد"
	UploadImageError        string = "خطا در آپلود تصویر"
	SomethingWrongHappened  string = "خطایی غیرمنتظره رخ داده است"
	SuccessfullyCreated     string = "باموفقیت ایجاد گردید"
	SuccessfullyUpdated     string = "با موفقیت بروزرسانی گردید"
	StoreImageOnDiskFailed  string = "خطا در ذخیره عکس بر روی هارددیسک"
	IDIsNotCorrect          string = "شناسه صحیح نمی باشد"
)

func HandleError(err error, notFoundMsg string) CustomError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(err.Error(), notFoundMsg, 404)
	}
	return New(err.Error(), InternalServerError, 500)
}
