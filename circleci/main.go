package circleci

import (
	"fmt"
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	"github.com/genofire/logmania/log"
	xmpp "github.com/mattn/go-xmpp"

	"github.com/genofire/hook2xmpp/config"
	ownXMPP "github.com/genofire/hook2xmpp/xmpp"
)

type Handler struct {
	client *xmpp.Client
	hooks  map[string]config.Hook
}

func NewHandler(client *xmpp.Client, newHooks []config.Hook) *Handler {
	hooks := make(map[string]config.Hook)

	for _, hook := range newHooks {
		if hook.Type == "circleci" {
			hooks[hook.URL] = hook
		}
	}
	return &Handler{
		client: client,
		hooks:  hooks,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	libHTTP.Read(r, &body)
	payload := body["payload"].(map[string]interface{})
	vcsURL, ok := payload["vcs_url"].(string)
	if !ok {
		log.Error(r.Body)
		http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
		return
	}

	hook, ok := h.hooks[vcsURL]
	if !ok {
		log.Errorf("No hook found for: '%s'", vcsURL)
		http.Error(w, fmt.Sprintf("no configuration for circleci for url %s", vcsURL), http.StatusNotFound)
		return
	}
	status := payload["status"].(string)
	buildNum := payload["build_num"].(float64)
	buildURL := payload["build_url"].(string)
	buildTime := payload["build_time_millis"].(float64)
	subject := payload["subject"].(string)
	msg := fmt.Sprintf("[%s] %s (%0.fs) - #%0.f: %s \n%s", vcsURL, status, buildTime/1000, buildNum, subject, buildURL)

	log.New().AddField("type", "circleci").Info(msg)
	ownXMPP.Notify(h.client, hook, msg)
}
