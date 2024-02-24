package logger

import (
	"log"
	"os"
)

var infoLogger = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime)
var errorLogger = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime)

func Fatalf(text string, v ...any) {
	errorLogger.Fatalf(text, v...)
}

func PrintfError(text string, v ...any) {
	errorLogger.Printf(text, v...)
}

func Printf(text string, v ...any) {
	infoLogger.Printf(text, v...)
}

func Println(text string) {
	infoLogger.Println(text)
}

func Panicf(text string, v ...any) {
	errorLogger.Panicf(text, v...)
}
