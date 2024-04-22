package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Просто логгер
var Log = logrus.New()
var file os.File

func Init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	Log.SetLevel(logrus.DebugLevel)

	if err != nil {
		panic(err)
	} else {
		Log.SetOutput(file)
	}
}

func End() {
	file.Close()
}
