package git

import (
	"fmt"

	"net/http"

	"github.com/bdlm/log"
	"github.com/mitchellh/mapstructure"
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

			var request requestBody
			if err := mapstructure.Decode(body, &request); err != nil {
				logger.Errorf("no readable payload: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithFields(map[string]interface{}{
				"url": request.Repository.HTMLURL,
				"msg": request.String(event),
			})

			ok := false
			for _, hook := range hooks {
				if request.Repository.HTMLURL != hook.URL {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, request.String(event))
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.Repository.HTMLURL), http.StatusNotFound)
			}
		}
	}
}
