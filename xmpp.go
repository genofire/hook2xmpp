package main

import (
	"github.com/bdlm/log"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

var client *xmpp.Client
var mucs []string

func notify(text string) {
	msg := stanza.Message{
		Attrs: stanza.Attrs{Type: stanza.MessageTypeGroupchat},
		Body:  text,
	}

	for _, muc := range config.StartupNotifyMuc {
		msg.To = muc
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"muc": muc,
				"msg": text,
			}).Errorf("error on startup notify: %s", err)
		}
	}

	msg.Type = stanza.MessageTypeChat
	for _, user := range config.StartupNotifyUser {
		msg.To = user
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"user": user,
				"msg":  text,
			}).Errorf("error on startup notify: %s", err)
		}
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

	return client.Send(stanza.Presence{Attrs: stanza.Attrs{To: jid},
		Extensions: []stanza.PresExtension{
			stanza.MucPresence{
				History: stanza.History{MaxStanzas: stanza.NewNullableInt(0)},
			}},
	})

}

func postStartup(c xmpp.StreamClient) {
	for _, muc := range config.StartupNotifyMuc {
		if err := joinMUC(muc, config.Nickname); err != nil {
			log.WithField("muc", muc).Errorf("error on joining muc: %s", err)
		}
	}
	for _, hooks := range config.Hooks {
		for _, hook := range hooks {
			for _, muc := range hook.NotifyMuc {
				if err := joinMUC(muc, config.Nickname); err != nil {
					log.WithField("muc", muc).Errorf("error on joining muc: %s", err)
				}
			}
		}
	}
	notify("started hock2xmpp")
}

func closeXMPP() {
	notify("stopped of hock2xmpp")

	for _, muc := range mucs {
		if err := client.Send(stanza.Presence{Attrs: stanza.Attrs{
			To:   muc,
			Type: stanza.PresenceTypeUnavailable,
		}}); err != nil {
			log.WithField("muc", muc).Errorf("error on leaving muc: %s", err)
		}
	}

}
