package main

import (
	"github.com/bdlm/log"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

var mucs []string

func notify(c xmpp.Sender, text string) {
	msg := stanza.Message{
		Attrs: stanza.Attrs{Type: stanza.MessageTypeGroupchat},
		Body:  text,
	}

	for _, muc := range config.StartupNotifyMuc {
		msg.To = muc
		if err := c.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"muc": muc,
				"msg": text,
			}).Errorf("error on startup notify: %s", err)
		}
	}

	msg.Type = stanza.MessageTypeChat
	for _, user := range config.StartupNotifyUser {
		msg.To = user
		if err := c.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"user": user,
				"msg":  text,
			}).Errorf("error on startup notify: %s", err)
		}
	}
	log.Infof("notify: %s", text)
}

func joinMUC(c xmpp.Sender, to, nick string) error {

	toJID, err := xmpp.NewJid(to)
	if err != nil {
		return err
	}
	toJID.Resource = nick
	jid := toJID.Full()

	mucs = append(mucs, jid)

	return c.Send(stanza.Presence{Attrs: stanza.Attrs{To: jid},
		Extensions: []stanza.PresExtension{
			stanza.MucPresence{
				History: stanza.History{MaxStanzas: stanza.NewNullableInt(0)},
			}},
	})

}

func postStartup(c xmpp.Sender) {
	for _, muc := range config.StartupNotifyMuc {
		if err := joinMUC(c, muc, config.Nickname); err != nil {
			log.WithField("muc", muc).Errorf("error on joining muc: %s", err)
		}
	}
	for _, hooks := range config.Hooks {
		for _, hook := range hooks {
			for _, muc := range hook.NotifyMuc {
				if err := joinMUC(c, muc, config.Nickname); err != nil {
					log.WithField("muc", muc).Errorf("error on joining muc: %s", err)
				}
			}
		}
	}
	notify(c, "started hock2xmpp")
}

func closeXMPP(c xmpp.Sender) {
	notify(c, "stopped of hock2xmpp")

	for _, muc := range mucs {
		if err := c.Send(stanza.Presence{Attrs: stanza.Attrs{
			To:   muc,
			Type: stanza.PresenceTypeUnavailable,
		}}); err != nil {
			log.WithField("muc", muc).Errorf("error on leaving muc: %s", err)
		}
	}
}
