package util

import (
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/go-playground/log/handlers/syslog"
	"github.com/lpimem/hlcsrv/conf"
)

func init() {
	if conf.IsDebug() {
		cLog := console.New()
		log.RegisterHandler(cLog, log.AllLevels...)
		return
	}

	sLog, err := syslog.New("", "", "hlcsrv", nil)
	if err != nil {
		panic(err)
	}
	levels := []log.Level{log.WarnLevel, log.ErrorLevel, log.AlertLevel, log.FatalLevel}
	log.RegisterHandler(sLog, levels...)
}
