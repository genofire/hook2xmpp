package runtime

type Config struct {
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

	StartupNotifyUser []string `toml:"startup_notify_user"`
	StartupNotifyMuc  []string `toml:"startup_notify_muc"`

	Hooks map[string][]Hook `toml:"hooks"`
}

type Hook struct {
	URL        string   `toml:"url"`
	NotifyUser []string `toml:"notify_user"`
	NotifyMuc  []string `toml:"notify_muc"`
}
