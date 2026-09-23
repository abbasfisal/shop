package logging

var GlobalLog Logger

type Logger interface {
	Init()
	Debug(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	DebugF(template string, args ...interface{})
	Info(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	InfoF(template string, args ...interface{})
	Warn(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	WarnF(template string, args ...interface{})
	Error(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	ErrorF(template string, args ...interface{})
	Fatal(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	FatalF(template string, args ...interface{})
}

func NewLogger() Logger {
	return NewZapLogger()
}
