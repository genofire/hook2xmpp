package xmpp

import (
	"log"

	"github.com/genofire/hook2xmpp/config"

	xmpp "github.com/mattn/go-xmpp"
)

func Start(client *xmpp.Client) {
	for {
		m, err := client.Recv()
		if err != nil {
			continue
		}
		switch v := m.(type) {
		case xmpp.Chat:
			if v.Type == "chat" {
				log.Printf("from %s: %s", v.Remote, v.Text)
			}
			if v.Type == "groupchat" {
			}
		case xmpp.Presence:
			// do nothing
		}
	}
}

func Notify(client *xmpp.Client, hook config.Hook, msg string) {
	for _, muc := range hook.NotifyMuc {
		client.SendHtml(xmpp.Chat{Remote: muc, Type: "groupchat", Text: msg})
	}
	for _, user := range hook.NotifyUser {
		client.SendHtml(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
	}
}
