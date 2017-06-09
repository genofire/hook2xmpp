package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"github.com/genofire/golang-lib/log"
)

type Config struct {
	WebserverBind string `toml:"webserver_bind"`

	XMPP struct {
		Host          string `toml:"host"`
		Username      string `toml:"username"`
		Password      string `toml:"password"`
		Debug         bool   `toml:"debug"`
		StartupNotify string `toml:"startup_notify"`
	} `toml:"xmpp"`

	Hooks []Hook `toml:"hooks"`
}

type Hook struct {
	NotifyUser []string `toml:"notify_user"`
	NotifyMuc  []string `toml:"notify_muc"`

	Type   string `toml:"type"`
	Github struct {
		Project string `toml:"project"`
	} `toml:"github"`
}

func ReadConfigFile(path string) *Config {
	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Log.Panic(err)
	}
	if err := toml.Unmarshal(file, config); err != nil {
		log.Log.Panic(err)
	}

	return config
}
