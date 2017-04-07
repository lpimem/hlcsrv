package util

import (
	"log"
	"os"
)

var _logger *log.Logger

func Log(v ...interface{}) {
	_logger.Println(v...)
}

func Debug(v ...interface{}) {
	_logger.Println(v...)
}

func init() {
	_logger = log.New(os.Stdout, "hlcsrv: ", log.Ldate|log.Ltime|log.Lshortfile)
}
