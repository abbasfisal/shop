package custom_error

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CustomError struct {
	OriginalMessage string
	DisplayMessage  string //will show to user
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
	UpdateWasFailed         string = "خطا در بروزرسانی"
	StoreImageOnDiskFailed  string = "خطا در ذخیره عکس بر روی هارددیسک"
	IDIsNotCorrect          string = "شناسه صحیح نمی باشد"
	IRMobileIsInvalid              = "شماره موبایل معتبر نیست"
	//OTP

	OTPRequestTooSoon = "باید ۴ دقیقه صبر کنید برای دریافت کد جدید"
	OTPTooManyRequest = "شما در حال حاضر به سقف درخواست‌ها برای این شماره در یک ساعت گذشته رسیده‌اید"
)
const (
	OTPTooSoonCode            = 4152
	OTPTooManyRequestCode     = 4150
	OtpSomethingGoesWrongCode = 4053

	//sesion

	CreateSessionFailedCode = 8001
)

func HandleError(err error, notFoundMsg string) CustomError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(err.Error(), notFoundMsg, 404)
	}
	return New(err.Error(), InternalServerError, 500)
}

// DuplicateProductInventory این خطا
// زمان ایجاد موجودی ایجاد میشه در صورتی که کاربر هیچگونه اتریبیوت-ویژگی رو پاس نداده باشه
// و صرفا بخواد که فقط موجودی برای محصول اضافه کنه
type DuplicateProductInventory struct {
	ProductID uint
}

func (d *DuplicateProductInventory) Error() string {
	return fmt.Sprintf("duplicate product inventory found with ID: %d", d.ProductID)
}
