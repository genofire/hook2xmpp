package git

import (
	"fmt"

	"net/http"

	libHTTP "dev.sum7.eu/genofire/golang-lib/http"
	"github.com/bdlm/log"
	"github.com/mitchellh/mapstructure"
	"gosrc.io/xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

var eventHeader = map[string]string{
	"X-GitHub-Event": "X-Hub-Signature",
	"X-Gogs-Event":   "X-Gogs-Delivery",
}

const hookType = "git"

func init() {
	runtime.HookRegister[hookType] = func(client xmpp.Sender, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		log.WithField("type", hookType).Info("loaded")
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

			event := ""
			secret := ""
			for head, headSecret := range eventHeader {
				event = r.Header.Get(head)

				if event != "" {
					secret = r.Header.Get(headSecret)
					break
				}
			}

			var body map[string]interface{}
			libHTTP.Read(r, &body)

			if s, ok := body["secret"]; ok && secret == "" {
				secret = s.(string)
			}

			if event == "" || secret == "" {
				logger.Warnf("no secret or event found")
				http.Error(w, fmt.Sprintf("no secret or event found"), http.StatusNotFound)
				return
			}

			var request requestBody
			if err := mapstructure.Decode(body, &request); err != nil {
				logger.Errorf("no readable payload: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithFields(map[string]interface{}{
				"url": request.Repository.URL,
				"msg": request.String(event),
			})

			ok := false
			msg := request.String(event)

			for _, hook := range hooks {
				if secret != hook.Secret {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, msg, msg)
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.Repository.URL), http.StatusNotFound)
			}
		}
	}
}
