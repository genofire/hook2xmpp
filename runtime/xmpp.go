package runtime

import (
	"github.com/bdlm/log"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

func NotifyImage(client *xmpp.Client, hook Hook, url string, desc string) {
	msg := stanza.Message{
		Attrs: stanza.Attrs{Type: stanza.MessageTypeGroupchat},
		Body:  url,
		Extensions: []stanza.MsgExtension{
			stanza.OOB{URL: url, Desc: desc},
		},
	}

	for _, muc := range hook.NotifyMuc {
		msg.To = muc
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"muc": muc,
				"url": url,
			}).Errorf("error on image notify: %s", err)
		}
	}

	msg.Type = stanza.MessageTypeChat
	for _, user := range hook.NotifyUser {
		msg.To = user
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"user": user,
				"url":  url,
			}).Errorf("error on image notify: %s", err)
		}
	}
}

func Notify(client *xmpp.Client, hook Hook, text, html string) {
	msg := stanza.Message{
		Attrs: stanza.Attrs{Type: stanza.MessageTypeGroupchat},
		Body:  text,
		Extensions: []stanza.MsgExtension{
			stanza.HTML{Body: stanza.HTMLBody{InnerXML: html}},
		},
	}

	for _, muc := range hook.NotifyMuc {
		msg.To = muc
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"muc":  muc,
				"text": text,
			}).Errorf("error on notify: %s", err)
		}
	}

	msg.Type = stanza.MessageTypeChat
	for _, user := range hook.NotifyUser {
		msg.To = user
		if err := client.Send(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"user": user,
				"text": text,
			}).Errorf("error on notify: %s", err)
		}
	}
}
