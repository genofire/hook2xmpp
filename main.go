package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/golang-lib/file"

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
		log.Panic(err)
	}

	log.Printf("Started hock2xmpp with %s", client.JID())

	client.SendHtml(xmpp.Chat{Remote: config.XMPP.StartupNotify, Type: "chat", Text: "startup of hock2xmpp"})
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

	// Wait for system signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	client.SendHtml(xmpp.Chat{Remote: config.XMPP.StartupNotify, Type: "chat", Text: "stopped of hock2xmpp"})

	srv.Close()

	log.Print("received", sig)
}
