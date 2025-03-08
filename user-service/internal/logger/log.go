package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func Init() {
	logFile, err := os.OpenFile("post_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Log faylini ochishda xato:", err)
	}

	Logger = log.New(logFile, "POSTGRES: ", log.Ldate|log.Ltime|log.Lshortfile)
}