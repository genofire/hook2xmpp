package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
	"gosrc.io/xmpp"

	_ "dev.sum7.eu/genofire/hook2xmpp/circleci"
	_ "dev.sum7.eu/genofire/hook2xmpp/git"
	_ "dev.sum7.eu/genofire/hook2xmpp/gitlab"
	_ "dev.sum7.eu/genofire/hook2xmpp/grafana"
	_ "dev.sum7.eu/genofire/hook2xmpp/prometheus"
	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

var config = runtime.Config{}

func main() {
	configFile := "config.conf"
	flag.StringVar(&configFile, "config", configFile, "path of configuration file")
	flag.Parse()

	if err := file.ReadTOML(configFile, &config); err != nil {
		log.WithField("tip", "maybe call me with: hook2xmpp--config /etc/hook2xmpp.conf").Panicf("error on read config: %s", err)
	}

	log.SetLevel(config.LogLevel)

	router := xmpp.NewRouter()
	client, err := xmpp.NewClient(xmpp.Config{
		Address:  config.XMPP.Address,
		Jid:      config.XMPP.JID,
		Password: config.XMPP.Password,
	}, router)

	if err != nil {
		log.Panicf("error on startup xmpp client: %s", err)
	}

	cm := xmpp.NewStreamManager(client, postStartup)
	go func() {
		err := cm.Run()
		log.Panic("closed connection:", err)
	}()
	for hookType, getHandler := range runtime.HookRegister {
		hooks, ok := config.Hooks[hookType]
		if ok {
			http.HandleFunc("/"+hookType, getHandler(client, hooks))
		}
	}

	srv := &http.Server{
		Addr: config.WebserverBind,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for system signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	closeXMPP()

	srv.Close()

	log.Infof("closed by receiving: %s", sig)
}
