package circleci

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bdlm/log"
	libHTTP "github.com/genofire/golang-lib/http"
	xmpp "github.com/mattn/go-xmpp"
	"github.com/mitchellh/mapstructure"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
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

func init() {
	runtime.HookRegister[hookType] = func(client *xmpp.Client, hooks []runtime.Hook) func(w http.ResponseWriter, r *http.Request) {
		log.WithField("type", hookType).Info("loaded")
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("type", hookType)

			var body interface{}
			libHTTP.Read(r, &body)

			var request requestBody
			if err := mapstructure.Decode(body, &request); err != nil {
				logger.Errorf("no readable payload: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}
			logger = logger.WithFields(map[string]interface{}{
				"url": request.RuleURL,
				"msg": request.String(),
			})

			ruleURL, err := url.Parse(request.RuleURL)
			if err != nil {
				logger.Errorf("could not parse ruleURL: %s", err)
				http.Error(w, fmt.Sprintf("no readable payload"), http.StatusInternalServerError)
				return
			}

			ok := false
			for _, hook := range hooks {
				if ruleURL.Hostname() != hook.URL {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, request.String())
				if request.ImageURL != "" {
				    runtime.NotifyImage(client, hook, request.ImageURL)
				}
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.RuleURL), http.StatusNotFound)
			}
		}
	}
}
