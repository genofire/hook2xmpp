package runtime

import (
	"net/http"

	xmpp "github.com/mattn/go-xmpp"
)

type HookHandler func(*xmpp.Client, []Hook) func(http.ResponseWriter, *http.Request)

var HookRegister map[string]HookHandler

func init() {
	HookRegister = make(map[string]HookHandler)
}
