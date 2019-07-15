package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	libHTTP "dev.sum7.eu/genofire/golang-lib/http"
	"github.com/bdlm/log"
	"gosrc.io/xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

var eventHeader = map[string]string{
	"X-GitHub-Event": "X-Hub-Signature",
	"X-Gogs-Event":   "X-Gogs-Delivery",
	"X-Gitlab-Event": "X-Gitlab-Token",
}

const hookType = "gitlab"

func init() {
	runtime.HookRegister[hookType] = func(client xmpp.Sender, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		log.WithField("type", hookType).Info("loaded")
		return func(w http.ResponseWriter, r *http.Request) {
			event := r.Header.Get("X-Gitlab-Event")
			secret := r.Header.Get("X-Gitlab-Token")

			logger := log.WithFields(map[string]interface{}{
				"event": event,
				"type":  hookType,
			})

			gitLabEvent := Event(event)

			if gitLabEvent == "" || secret == "" {
				logger.Warnf("no secret or event found")
				http.Error(w, fmt.Sprintf("no secret or event found"), http.StatusNotFound)
				return
			}

			var msg string
			var err error

			switch gitLabEvent {
			case PushEvents:
				var pl PushEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case TagEvents:
				var pl TagEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case ConfidentialIssuesEvents:
				var pl ConfidentialIssueEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case IssuesEvents:
				var pl IssueEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case CommentEvents:
				var pl CommentEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case MergeRequestEvents:
				var pl MergeRequestEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case WikiPageEvents:
				var pl WikiPageEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case PipelineEvents:
				var pl PipelineEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case BuildEvents:
				var pl BuildEventPayload
				err = libHTTP.Read(r, &pl)
				msg = pl.String()

			case SystemEvents:
				var data map[string]interface{}
				var buf bytes.Buffer
				tee := io.TeeReader(r.Body, &buf)
				if err = json.NewDecoder(tee).Decode(&data); err != nil {
					msg = fmt.Sprintf("unable to decode gitlab system event")
				} else if event, ok := data["event_name"]; ok {
					switch event {
					case "push":
						var pl PushEventPayload
						err = json.NewDecoder(&buf).Decode(&pl)
						msg = fmt.Sprintf("[S]%s", pl.String())
					default:
						err = nil
						msg = fmt.Sprintf("unknown gitlab system event '%s' received", event)
					}
				} else {
					err = nil
					msg = fmt.Sprintf("unable to get 'event_name' of gitlab '%s'", gitLabEvent)
				}

			default:
				err = nil
				msg = fmt.Sprintf("unknown gitlab event '%s' received", gitLabEvent)
			}

			logger = logger.WithField("msg", msg)

			if err != nil {
				logger.Warnf("unable decode message: %s", err)
				return
			}

			ok := false
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
				http.Error(w, fmt.Sprintf("no configuration for %s for message: %s", hookType, msg), http.StatusNotFound)
			}
		}
	}
}
