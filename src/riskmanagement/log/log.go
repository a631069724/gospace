package log

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLog(path string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	logger.Printf(format, v)
}
