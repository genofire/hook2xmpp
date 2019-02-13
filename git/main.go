package git

import (
	"fmt"

	"net/http"

	"github.com/bdlm/log"
	libHTTP "github.com/genofire/golang-lib/http"
	xmpp "github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

var eventHeader = []string{"X-GitHub-Event", "X-Gogs-Event"}

const hookType = "git"

func init() {
	runtime.HookRegister[hookType] = func(client *xmpp.Client, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

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

			var body map[string]interface{}
			libHTTP.Read(r, &body)

			repository := body["repository"].(map[string]interface{})
			url, ok := repository["html_url"].(string)
			if !ok {
				logger.Error("no readable payload")
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithField("url", url)

			msg := PayloadToString(event, body)
			logger = logger.WithField("msg", msg)

			ok = false
			for _, hook := range hooks {
				if url != hook.URL {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, msg)
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for git for url: %s", url), http.StatusNotFound)
			}
		}
	}
}
