package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/genofire/golang-lib/log"
	"github.com/mattn/go-xmpp"
	"github.com/pierrre/githubhook"

	"github.com/genofire/hook2xmpp/circleci"
	configuration "github.com/genofire/hook2xmpp/config"
	"github.com/genofire/hook2xmpp/github"
	ownXMPP "github.com/genofire/hook2xmpp/xmpp"
)

func main() {
	configFile := "config.conf"
	flag.StringVar(&configFile, "config", configFile, "path of configuration file")
	flag.Parse()

	// load config
	config := configuration.ReadConfigFile(configFile)
	client, err := xmpp.NewClientNoTLS(config.XMPP.Host, config.XMPP.Username, config.XMPP.Password, config.XMPP.Debug)
	if err != nil {
		log.Log.Panic(err)
	}

	log.Log.Infof("Started hock2xmpp with %s", client.JID())

	client.SendHtml(xmpp.Chat{Remote: config.XMPP.StartupNotify, Type: "chat", Text: "startup of hock2xmpp"})
	go ownXMPP.Start(client)

	githubHandler := github.NewHandler(client, config.Hooks)
	handler := &githubhook.Handler{
		Delivery: githubHandler.Deliviery,
	}
	http.Handle("/github", handler)
	circleciHandler := circleci.NewHandler(client, config.Hooks)
	http.Handle("/circleci", circleciHandler)

	srv := &http.Server{
		Addr: config.WebserverBind,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// Wait for system signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	client.SendHtml(xmpp.Chat{Remote: config.XMPP.StartupNotify, Type: "chat", Text: "stopped of hock2xmpp"})

	srv.Close()

	log.Log.Info("received", sig)
}
