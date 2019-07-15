package prometheus

import (
	"fmt"
	"strings"

	"net/http"

	libHTTP "dev.sum7.eu/genofire/golang-lib/http"
	"github.com/bdlm/log"
	"github.com/prometheus/alertmanager/notify/webhook"
	"gosrc.io/xmpp"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

const hookType = "prometheus"

func init() {
	runtime.HookRegister[hookType] = func(client xmpp.Sender, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		log.WithField("type", hookType).Info("loaded")
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

			var request webhook.Message
			if err := libHTTP.Read(r, &request); err != nil {
				logger.Errorf("no readable payload: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}

			// title
			content := strings.Join(request.GroupLabels.Values(), " ")
			html := fmt.Sprintf(`<span style="font-weight: bold;">%s</span>`, content)

			statusColor := "#ffff00"

			switch request.Status {
			case "resolved":
				statusColor = "#00ff00"
			case "firing":
				statusColor = "#ff8700"
			}

			firingAlerts := request.Alerts.Firing()
			if len(firingAlerts) > 0 {
				for _, a := range firingAlerts {
					if description, ok := a.Annotations["message"]; ok {
						content = fmt.Sprintf("%s\n%s", content, description)
						html = fmt.Sprintf("%s<br/>%s", html, description)
					}
				}
				content = fmt.Sprintf("[%s:%d] %s", request.Status, len(firingAlerts), content)
				html = fmt.Sprintf(`<span style="color:%s">%s:%d</span> %s`, statusColor, request.Status, len(firingAlerts), html)
			} else {
				html = fmt.Sprintf(`<span style="color:%s">%s</span> %s`, statusColor, request.Status, content)
				content = fmt.Sprintf("[%s] %s", request.Status, content)
			}

			logger = logger.WithField("body", content)

			ok := false
			token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
			for _, hook := range hooks {
				if token != hook.Secret {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, content, html)
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.ExternalURL), http.StatusNotFound)
			}
		}
	}
}
