package runtime

import "github.com/bdlm/std/logger"

type Config struct {
	LogLevel      logger.Level `toml:"log_level"`
	WebserverBind string       `toml:"webserver_bind"`

	XMPP struct {
		Host          string `toml:"host"`
		Username      string `toml:"username"`
		Resource      string `toml:"resource"`
		Password      string `toml:"password"`
		Debug         bool   `toml:"debug"`
		NoTLS         bool   `toml:"no_tls"`
		StartTLS      bool   `toml:"start_tls"`
		Session       bool   `toml:"session"`
		Status        string `toml:"status"`
		StatusMessage string `toml:"status_message"`
	} `toml:"xmpp"`

	Nickname string `toml:"nickname"`

	StartupNotifyUser []string `toml:"startup_notify_user"`
	StartupNotifyMuc  []string `toml:"startup_notify_muc"`

	Hooks map[string][]Hook `toml:"hooks"`
}

type Hook struct {
	Secret     string   `toml:"secret"`
	NotifyUser []string `toml:"notify_user"`
	NotifyMuc  []string `toml:"notify_muc"`
}
