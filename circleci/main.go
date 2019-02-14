package circleci

import (
	"fmt"

	"net/http"

	libHTTP "dev.sum7.eu/genofire/golang-lib/http"
	"github.com/bdlm/log"
	xmpp "github.com/mattn/go-xmpp"
	"github.com/mitchellh/mapstructure"

	"dev.sum7.eu/genofire/hook2xmpp/runtime"
)

const hookType = "circleci"

type requestBody struct {
	Payload struct {
		VCSURL    string  `mapstructure:"vcs_url"`
		Status    string  `mapstructure:"status"`
		BuildNum  float64 `mapstructure:"build_num"`
		BuildURL  string  `mapstructure:"build_url"`
		BuildTime float64 `mapstructure:"build_time_millis"`
		Subject   string  `mapstructure:"subject"`
	} `mapstructure:"payload"`
}

func (r requestBody) String() string {
	return fmt.Sprintf("#%0.f (%0.fs): %s", r.Payload.BuildNum, r.Payload.BuildTime/1000, r.Payload.Subject)
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
				"url": request.Payload.VCSURL,
				"msg": request.String(),
			})

			ok := false
			for _, hook := range hooks {
				if request.Payload.VCSURL != hook.Secret {
					continue
				}
				logger.Infof("run hook")
				runtime.Notify(client, hook, request.String())
				ok = true
			}
			if !ok {
				logger.Warnf("no hook found")
				http.Error(w, fmt.Sprintf("no configuration for %s for url: %s", hookType, request.Payload.VCSURL), http.StatusNotFound)
			}
		}
	}
}
