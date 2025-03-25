package logging

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
	Log  *logrus.Logger
)

func InitLogrus() *logrus.Logger {
	once.Do(func() {

		w, _ := os.Getwd()                                          //get working director
		rootPath, err := filepath.Abs(filepath.Join(w, "..", "..")) // get absolute path
		if err != nil {
			log.Fatal("abs failed:", err)
		}

		// store log path -> your-root-project/storage/logs/shop.log
		logPath := filepath.Join(rootPath, "storage", "logs", "shop.log")

		fileWriter := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10, //MB
			MaxAge:     10, //day
			MaxBackups: 5,
			LocalTime:  false,
			Compress:   true,
		}
		multi := io.MultiWriter(os.Stdout, fileWriter)

		Log = logrus.New()
		Log.SetFormatter(&logrus.JSONFormatter{DataKey: "MetaData"})
		Log.SetOutput(multi)
	})
	return Log
}
