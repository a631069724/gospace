package logger

import "log"

var MyLoger *log.Logger

func Println(v ...interface{}) {
	MyLoger.Println(v)
}
