package circleci

import (
	"fmt"

	"net/http"

	"github.com/bdlm/log"
	libHTTP "github.com/genofire/golang-lib/http"
	xmpp "github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

const hookType = "git"

func init() {
	runtime.HookRegister[hookType] = func(client *xmpp.Client, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

			var body map[string]interface{}
			libHTTP.Read(r, &body)

			payload := body["payload"].(map[string]interface{})
			url, ok := payload["vcs_url"].(string)
			if !ok {
				logger.Error("no readable payload")
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithField("url", url)

			status := payload["status"].(string)
			buildNum := payload["build_num"].(float64)
			buildURL := payload["build_url"].(string)
			buildTime := payload["build_time_millis"].(float64)
			subject := payload["subject"].(string)
			msg := fmt.Sprintf("[%s] %s (%0.fs) - #%0.f: %s \n%s", url, status, buildTime/1000, buildNum, subject, buildURL)
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
				http.Error(w, fmt.Sprintf("no configuration for circleci for url: %s", url), http.StatusNotFound)
			}
		}
	}
}
