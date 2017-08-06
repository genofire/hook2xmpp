package config

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Syslog struct {
		Enable  bool   `toml:"enable"`
		Type    string `toml:"type"`
		Address string `toml:"address"`
	} `toml:"syslog"`

	WebserverBind string `toml:"webserver_bind"`

	XMPP struct {
		Host          string `toml:"host"`
		Username      string `toml:"username"`
		Password      string `toml:"password"`
		Debug         bool   `toml:"debug"`
		NoTLS         bool   `toml:"no_tls"`
		Session       bool   `toml:"session"`
		Status        string `toml:"status"`
		StatusMessage string `toml:"status_message"`
		StartupNotify string `toml:"startup_notify"`
	} `toml:"xmpp"`

	Hooks []Hook `toml:"hooks"`
}

type Hook struct {
	Type       string   `toml:"type"`
	URL        string   `toml:"url"`
	NotifyUser []string `toml:"notify_user"`
	NotifyMuc  []string `toml:"notify_muc"`
}

func ReadConfigFile(path string) *Config {
	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	if err := toml.Unmarshal(file, config); err != nil {
		log.Panic(err)
	}

	return config
}
