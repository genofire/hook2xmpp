package runtime

import (
	"net/http"

	"gosrc.io/xmpp"
)

type HookHandler func(*xmpp.Client, []Hook) func(http.ResponseWriter, *http.Request)

var HookRegister map[string]HookHandler

func init() {
	HookRegister = make(map[string]HookHandler)
}
