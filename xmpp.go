package main

import (
	"github.com/bdlm/log"
	"gosrc.io/xmpp"
)

var client *xmpp.Client
var mucs []string

func notify(text string) {
	msg := xmpp.Message{
		Attrs: xmpp.Attrs{Type: xmpp.MessageTypeGroupchat},
		Body:  text,
	}

	for _, muc := range config.StartupNotifyMuc {
		msg.To = muc
		client.Send(msg)
	}

	msg.Type = xmpp.MessageTypeChat
	for _, user := range config.StartupNotifyUser {
		msg.To = user
		client.Send(msg)
	}
	log.Infof("notify: %s", text)
}

func joinMUC(to, nick string) error {

	toJID, err := xmpp.NewJid(to)
	if err != nil {
		return err
	}
	toJID.Resource = nick
	jid := toJID.Full()

	mucs = append(mucs, jid)

	return client.Send(xmpp.Presence{Attrs: xmpp.Attrs{To: jid},
		Extensions: []xmpp.PresExtension{
			xmpp.MucPresence{
				History: xmpp.History{MaxStanzas: xmpp.NewNullableInt(0)},
			}},
	})

}

func postStartup(c xmpp.StreamClient) {
	for _, muc := range config.StartupNotifyMuc {
		joinMUC(muc, config.Nickname)
	}
	for _, hooks := range config.Hooks {
		for _, hook := range hooks {
			for _, muc := range hook.NotifyMuc {
				joinMUC(muc, config.Nickname)
			}
		}
	}
	notify("started hock2xmpp")
}

func closeXMPP() {
	notify("stopped of hock2xmpp")

	for _, muc := range mucs {
		client.Send(xmpp.Presence{Attrs: xmpp.Attrs{
			To:   muc,
			Type: xmpp.PresenceTypeUnavailable,
		}})
	}

}
