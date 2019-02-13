package git

import (
	"fmt"
	"log"
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	xmpp "github.com/mattn/go-xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

var eventHeader = []string{"X-GitHub-Event", "X-Gogs-Event"}

func init() {
	runtime.HookRegister["git"] = func(client *xmpp.Client, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
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

			ok := false
			for _, hook := range hooks {
				if url != hook.URL {
					continue
				}
				log.Printf("run hook for git: %s", msg)
				runtime.Notify(client, hook, msg)
				ok = true
			}
			if !ok {
				log.Fatalf("No hook found for: '%s'", url)
				http.Error(w, fmt.Sprintf("no configuration for git for url %s", url), http.StatusNotFound)
			}
		}
	}
}
