package logger

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
}

func Println(msg string) {
	InfoLogger.Output(2, msg)
}

func Info(msg string) {
	InfoLogger.Output(2, msg)
}

func Warn(msg string) {
	WarningLogger.Output(2, msg)
}

func Error(err error) {
	ErrorLogger.Output(2, err.Error())
}

func Fatal(msg string) {
	ErrorLogger.Fatal(msg)
}
