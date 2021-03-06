package grafana

import (
	"fmt"
	"net/http"

	libHTTP "dev.sum7.eu/genofire/golang-lib/http"
	"github.com/bdlm/log"
	"gosrc.io/xmpp"

	"dev.sum7.eu/sum7/hook2xmpp/runtime"
)

const hookType = "grafana"

type evalMatch struct {
	Tags   map[string]string `mapstructure:"tags,omitempty"`
	Metric string            `mapstructure:"metric"`
	Value  float64           `mapstructure:"value"`
}

type requestBody struct {
	Title       string      `mapstructure:"title"`
	State       string      `mapstructure:"state"`
	RuleID      int64       `mapstructure:"ruleId"`
	RuleName    string      `mapstructure:"ruleName"`
	RuleURL     string      `mapstructure:"ruleUrl"`
	EvalMatches []evalMatch `mapstructure:"evalMatches"`
	ImageURL    string      `mapstructure:"imageUrl"`
	Message     string      `mapstructure:"message"`
}

func (r requestBody) String() string {
	msg := fmt.Sprintf("%s: %s", r.Title, r.Message)
	for _, e := range r.EvalMatches {
		msg = fmt.Sprintf("%s %s=%f", msg, e.Metric, e.Value)
	}
	return msg
}
func (r requestBody) HTML() string {
	stateColor := "#ffff00"
	switch r.State {
	case "alerting":
		stateColor = "#ff0000"
	case "ok":
		stateColor = "#00ff00"
	}

	msg := fmt.Sprintf("<span style=\"color: %s;\">%s</span> <span style=\"font-weight: bold;\">%s</span>: %s<br/>", stateColor, r.State, r.RuleName, r.Message)
	for _, e := range r.EvalMatches {
		msg = fmt.Sprintf("%s %s=%f", msg, e.Metric, e.Value)
	}
	return msg
}

func init() {
	runtime.HookRegister[hookType] = func(client xmpp.Sender, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		log.WithField("type", hookType).Info("loaded")
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

			_, secret, ok := r.BasicAuth()

			if !ok {
				logger.Errorf("no secret found")
				http.Error(w, fmt.Sprintf("no secret found (basic-auth password)"), http.StatusUnauthorized)
				return
			}

			var request requestBody
			if err := libHTTP.Read(r, &request); err != nil {
				logger.Errorf("no readable payload: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithFields(map[string]interface{}{
				"url":   request.RuleURL,
				"msg":   request.String(),
				"image": request.ImageURL,
			})

			ok = false
			msg := request.String()
			html := request.HTML()

			for _, hook := range hooks {
				if secret != hook.Secret {
					continue
				}

				runtime.Notify(client, hook, msg, html)
				if request.ImageURL != "" {
					runtime.NotifyImage(client, hook, request.ImageURL, request.String())
				} else {
				}
				logger.Infof("run hook")
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.RuleURL), http.StatusNotFound)
			}
		}
	}
}
