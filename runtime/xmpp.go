package runtime

import (
	"gosrc.io/xmpp"
)

func NotifyImage(client *xmpp.Client, hook Hook, url string, desc string) {
	msg := xmpp.Message{
		Attrs: xmpp.Attrs{Type: xmpp.MessageTypeGroupchat},
		Body:  url,
		Extensions: []xmpp.MsgExtension{
			xmpp.OOB{URL: url, Desc: desc},
		},
	}

	for _, muc := range hook.NotifyMuc {
		msg.To = muc
		client.Send(msg)
	}

	msg.Type = xmpp.MessageTypeChat
	for _, user := range hook.NotifyUser {
		msg.To = user
		client.Send(msg)
	}
}

func Notify(client *xmpp.Client, hook Hook, text, html string) {
	msg := xmpp.Message{
		Attrs: xmpp.Attrs{Type: xmpp.MessageTypeGroupchat},
		Body:  text,
		Extensions: []xmpp.MsgExtension{
			xmpp.HTML{Body: xmpp.HTMLBody{InnerXML: html}},
		},
	}

	for _, muc := range hook.NotifyMuc {
		msg.To = muc
		client.Send(msg)
	}

	msg.Type = xmpp.MessageTypeChat
	for _, user := range hook.NotifyUser {
		msg.To = user
		client.Send(msg)
	}
}
