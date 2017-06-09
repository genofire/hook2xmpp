package circleci

import (
	"fmt"
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	"github.com/genofire/golang-lib/log"
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
			repoFullName := fmt.Sprintf("%s/%s", hook.CircleCI.Username, hook.CircleCI.Reponame)
			hooks[repoFullName] = hook
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
	username, ok := payload["username"].(string)
	if !ok {
		log.Log.Error(r.Body)
		http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
		return
	}
	reponame := payload["reponame"].(string)
	repoFullName := fmt.Sprintf("%s/%s", username, reponame)

	hook, ok := h.hooks[repoFullName]
	if !ok {
		log.Log.Errorf("No hook found for: '%s'", repoFullName)
		http.Error(w, fmt.Sprintf("no configuration for circleci with username %s and reponame %s", username, reponame), http.StatusNotFound)
		return
	}
	status := payload["status"].(string)
	buildNum := payload["build_num"].(float64)
	buildURL := payload["build_url"].(string)
	buildTime := payload["build_time_millis"].(float64)
	subject := payload["subject"].(string)
	msg := fmt.Sprintf("[%s/%s] %s (%0.fms) - #%0.f: %s \n%s", username, reponame, status, buildTime, buildNum, subject, buildURL)

	log.Log.WithField("type", "circleci").Print(msg)
	ownXMPP.Notify(h.client, hook, msg)
}