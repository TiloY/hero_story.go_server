package log

import (
	"fmt"
	"log"
)

var writer *dailyFileWriter

var infoLogger, errorLogger *log.Logger

func Config() {
	writer = &dailyFileWriter{}
	infoLogger = log.New(writer, "[INFO]", 0)
}

func Info(format string, valArray ...interface{}) {
	_ = infoLogger.Output(
		2,
		fmt.Sprintf(format, valArray...),
	)
}

func Error(format string, valArray ...interface{}) {
	_ = errorLogger.Output(
		2,
		fmt.Sprintf(format, valArray...),
	)
}
