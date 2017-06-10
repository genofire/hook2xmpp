package git

import (
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	"github.com/genofire/golang-lib/log"
	xmpp "github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/config"
	ownXMPP "dev.sum7.eu/genofire/hook2xmpp/xmpp"
)

type Handler struct {
	client *xmpp.Client
	hooks  map[string]config.Hook
}

func NewHandler(client *xmpp.Client, newHooks []config.Hook) *Handler {
	hooks := make(map[string]config.Hook)

	for _, hook := range newHooks {
		if hook.Type == "git" {
			hooks[hook.URL] = hook
		}
	}
	return &Handler{
		client: client,
		hooks:  hooks,
	}
}

var eventHeader = []string{"X-GitHub-Event", "X-Gogs-Event"}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	event := ""
	for _, head := range eventHeader {
		event = r.Header.Get(head)

		if event != "" {
			break
		}
	}

	if event == "status" {
		return
	}

	libHTTP.Read(r, &payload)
	msg := PayloadToString(event, payload)
	repository := payload["repository"].(map[string]interface{})
	url := repository["html_url"].(string)

	hook, ok := h.hooks[url]
	if !ok {
		log.Log.Errorf("No hook found for: '%s'", url)
		return
	}

	log.Log.WithField("type", "git").Print(msg)
	ownXMPP.Notify(h.client, hook, msg)
}
