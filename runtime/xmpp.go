package runtime

import (
	"fmt"

	"github.com/bdlm/log"
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
				log.Debugf("from %s: %s", v.Remote, v.Text)
			}
			if v.Type == "groupchat" {
			}
		case xmpp.Presence:
			// do nothing
		}
	}
}
func NotifyImage(client *xmpp.Client, hook Hook, url string, desc string) {
	msg := fmt.Sprintf(`<message to='%%s' type='%%s'>
		<body>%s</body>
		<x xmlns='jabber:x:oob'>
			<url>%s</url>
			<desc>%s</desc>
		</x>
	</message>`, url, url, desc)

	for _, muc := range hook.NotifyMuc {
		client.SendOrg(fmt.Sprintf(msg, muc, "groupchat"))
	}
	for _, user := range hook.NotifyUser {
		client.SendOrg(fmt.Sprintf(msg, user, "chat"))
	}
}
func Notify(client *xmpp.Client, hook Hook, msg string) {
	for _, muc := range hook.NotifyMuc {
		client.SendHtml(xmpp.Chat{Remote: muc, Type: "groupchat", Text: msg})
	}
	for _, user := range hook.NotifyUser {
		client.SendHtml(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
	}
}
