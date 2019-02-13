package circleci

import (
	"fmt"
	"log"
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	xmpp "github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

func init() {
	runtime.HookRegister["circleci"] = func(client *xmpp.Client, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			libHTTP.Read(r, &body)
			payload := body["payload"].(map[string]interface{})
			vcsURL, ok := payload["vcs_url"].(string)
			if !ok {
				log.Fatal("no readable payload")
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}

			status := payload["status"].(string)
			buildNum := payload["build_num"].(float64)
			buildURL := payload["build_url"].(string)
			buildTime := payload["build_time_millis"].(float64)
			subject := payload["subject"].(string)
			msg := fmt.Sprintf("[%s] %s (%0.fs) - #%0.f: %s \n%s", vcsURL, status, buildTime/1000, buildNum, subject, buildURL)

			ok = false
			for _, hook := range hooks {
				if vcsURL != hook.URL {
					continue
				}
				log.Printf("run hook for circleci: %s", msg)
				runtime.Notify(client, hook, msg)
				ok = true
			}
			if !ok {
				log.Fatalf("No hook found for: '%s'", vcsURL)
				http.Error(w, fmt.Sprintf("no configuration for circleci for url %s", vcsURL), http.StatusNotFound)
			}
		}
	}
}
