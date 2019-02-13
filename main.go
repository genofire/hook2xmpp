package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
	"github.com/mattn/go-xmpp"

	_ "dev.sum7.eu/genofire/hook2xmpp/circleci"
	_ "dev.sum7.eu/genofire/hook2xmpp/git"
	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

func main() {
	configFile := "config.conf"
	flag.StringVar(&configFile, "config", configFile, "path of configuration file")
	flag.Parse()

	config := &runtime.Config{}

	if err := file.ReadTOML(configFile, config); err != nil {
		log.Panicf("error on read config: %s", err)
	}

	// load config
	options := xmpp.Options{
		Host:          config.XMPP.Host,
		User:          config.XMPP.Username,
		Password:      config.XMPP.Password,
		NoTLS:         config.XMPP.NoTLS,
		Debug:         config.XMPP.Debug,
		Session:       config.XMPP.Session,
		Status:        config.XMPP.Status,
		StatusMessage: config.XMPP.StatusMessage,
	}
	client, err := options.NewClient()
	if err != nil {
		log.Panicf("error on startup xmpp client: %s",err)
	}

	go runtime.Start(client)

	for hookType, getHandler := range runtime.HookRegister {
		hooks, ok := config.Hooks[hookType]
		if ok {
			http.HandleFunc(hookType, getHandler(client, hooks))
		}
	}

	srv := &http.Server{
		Addr: config.WebserverBind,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	var mucs []string
	for _, muc := range config.StartupNotifyMuc {
		mucs = append(mucs, muc)
		client.JoinMUCNoHistory(muc, config.Nickname)
	}
	for _, hooks := range config.Hooks {
		for _, hook := range hooks {
			for _, muc := range hook.NotifyMuc {
				mucs = append(mucs, muc)
				client.JoinMUCNoHistory(muc, config.Nickname)
			}
		}
	}

	notify := func (msg string) {
		for _, muc := range config.StartupNotifyMuc {
			client.SendHtml(xmpp.Chat{Remote: muc, Type: "groupchat", Text: msg})
		}
		for _, user := range config.StartupNotifyUser {
			client.SendHtml(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
		}
	}

	notify("startup of hock2xmpp")

	log.Infof("started hock2xmpp with %s", client.JID())

	// Wait for system signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	notify("stopped of hock2xmpp")
	
	for _, muc := range mucs {
		client.LeaveMUC(muc)
	}

	srv.Close()
	client.Close()

	log.Infof("closed by receiving: %s", sig)
}
