package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func InitLogger() {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Log = log.New(file, "", log.LstdFlags)
}
