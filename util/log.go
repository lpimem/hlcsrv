package util

import (
	"log"
	"os"
)

var _logger *log.Logger

const (
	_ = iota
	_DEBUG
	_INFO
	_WARN
	_ERROR
)

var _level = _DEBUG

func _log(lv int, v ...interface{}) {
	if lv >= _level {
		if lv >= _WARN {
			_logger.Panic(v...)
		} else {
			_logger.Println(v...)
		}
	}
}

func Debug(v ...interface{}) {
	_log(_DEBUG, v...)
}

func Log(v ...interface{}) {
	_log(_INFO, v...)
}

func Warn(v ...interface{}) {
	_log(_WARN, v...)
}

func Error(v ...interface{}) {
	_log(_ERROR, v...)
}

func init() {
	_logger = log.New(os.Stdout, "hlcsrv: ", log.Ldate|log.Ltime|log.Lshortfile)

	// log.SetOutput(os.Stdout)
	// log.SetPrefix("hlcsrv: ")
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
