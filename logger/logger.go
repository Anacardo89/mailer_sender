package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func CreateLogger() error {
	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access INFO log file:", err)
	}
	warnFile, err := os.OpenFile("logs/warn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access WARN log file:", err)
	}
	errorFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access ERROR log file:", err)
	}
	Info = log.New(infoFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(warnFile, "WARN:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
