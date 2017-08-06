package syslog

import (
	"log"
	"log/syslog"

	"github.com/genofire/hook2xmpp/config"
)

func Bind(config *config.Config) {
	sysLog, err := syslog.Dial(config.Syslog.Type, config.Syslog.Address, syslog.LOG_WARNING|syslog.LOG_DAEMON, "hook2xmpp")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(sysLog)
}
