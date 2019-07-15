package runtime

import (
	"net/http"

	"gosrc.io/xmpp"
)

type HookHandler func(xmpp.Sender, []Hook) func(http.ResponseWriter, *http.Request)

var HookRegister map[string]HookHandler

func init() {
	HookRegister = make(map[string]HookHandler)
}
