package util

import (
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/go-playground/log/handlers/syslog"
)

func init() {
	cLog := console.New()
	log.RegisterHandler(cLog, log.AllLevels...)

	sLog, err := syslog.New("", "", "hlcsrv", nil)
	if err != nil {
		panic(err)
	}
	log.RegisterHandler(sLog, log.WarnLevel, log.ErrorLevel, log.AlertLevel)
}
