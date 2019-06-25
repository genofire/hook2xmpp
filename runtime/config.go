package runtime

import "github.com/bdlm/std/logger"

type Config struct {
	LogLevel      logger.Level `toml:"log_level"`
	WebserverBind string       `toml:"webserver_bind"`

	XMPP struct {
		Address  string `toml:"address"`
		JID      string `toml:"jid"`
		Password string `toml:"password"`
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
