package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var zapSinLogger *zap.SugaredLogger

type ZapLogger struct {
	//config *config.Config
	logger *zap.SugaredLogger
}

var zapLogLevelMapping = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func NewZapLogger() *ZapLogger {
	logger := &ZapLogger{}
	logger.Init()
	return logger
}

func (l *ZapLogger) getLogLevel() zapcore.Level {
	level, exists := zapLogLevelMapping["info"]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func (l *ZapLogger) Init() {

	stdoutWriter := zapcore.Lock(os.Stdout)

	logDir := "./storage/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Create a file to write logs (if it doesn't exist, create it)
	logFile := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Setup the writers for stdout and file
	fileWriter := zapcore.AddSync(file)
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(stdoutWriter, fileWriter)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	stdoutCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		multiWriteSyncer,
		l.getLogLevel(),
	)
	logger := zap.New(stdoutCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	zapSinLogger = logger.With("AppName", "MyApp").With("LoggerName", "ZapLog")

	l.logger = zapSinLogger
}

func prepareLogKeys(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["category"] = cat
	extra["subCategory"] = sub
	return mapToZapParams(extra)
}

func mapToZapParams(keys map[ExtraKey]interface{}) []interface{} {
	params := make([]interface{}, 0, len(keys))
	for k, v := range keys {
		params = append(params, string(k), v)
	}
	return params
}

func (l *ZapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	l.logger.Debugw(msg, params...)
}

func (l *ZapLogger) DebugF(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *ZapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	l.logger.Infow(msg, params...)
}

func (l *ZapLogger) InfoF(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *ZapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	l.logger.Warnw(msg, params...)
}

func (l *ZapLogger) WarnF(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *ZapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	l.logger.Errorw(msg, params...)
}

func (l *ZapLogger) ErrorF(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *ZapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	l.logger.Fatalw(msg, params...)
}

func (l *ZapLogger) FatalF(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}
