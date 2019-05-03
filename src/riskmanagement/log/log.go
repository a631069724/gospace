package log

import (
	"log"
	"os"
)

var logger *log.Logger
var file *os.File

func InitLog(path string) {
	var err error
	file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	log.SetOutput(file)
}

func Info(format string, v ...interface{}) {
	logger.Printf(format, v)
}

func Close() {
	file.Close()
}
