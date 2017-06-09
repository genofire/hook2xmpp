package github

import (
	"fmt"
	"net/http"
	"strings"

	libHTTP "github.com/genofire/golang-lib/http"
	"github.com/genofire/golang-lib/log"
	xmpp "github.com/mattn/go-xmpp"

	"github.com/genofire/hook2xmpp/config"
	ownXMPP "github.com/genofire/hook2xmpp/xmpp"
)

type Handler struct {
	client *xmpp.Client
	hooks  map[string]config.Hook
}

func NewHandler(client *xmpp.Client, newHooks []config.Hook) *Handler {
	hooks := make(map[string]config.Hook)

	for _, hook := range newHooks {
		if hook.Type == "github" {
			hooks[hook.Github.Project] = hook
		}
	}
	return &Handler{
		client: client,
		hooks:  hooks,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	event := r.Header.Get("X-GitHub-Event")

	if event == "status" {
		return
	}

	libHTTP.Read(r, &payload)
	msg := PayloadToString(event, payload)
	repository := payload["repository"].(map[string]interface{})
	repoName := repository["full_name"].(string)

	hook, ok := h.hooks[repoName]
	if !ok {
		log.Log.Errorf("No hook found for: '%s'", repoName)
		return
	}

	log.Log.WithField("type", "github").Print(msg)
	ownXMPP.Notify(h.client, hook, msg)
}

var eventMsg = map[string]string{
	"commit_comment_created": "Commit comment",
	"status_error":           "Commit status: error",
	"status_failure":         "Commit status: failure",
	"status_pending":         "Commit status: pending",
	"status_success":         "Commit status: success",
	"create_branch":          "Create branch",
	"create_tag":             "Create tag",
	"delete_branch":          "Delete branch",
	"delete_tag":             "Delete tag",
	"issue_comment_created":  "Issue comment",
	"issue_comment_deleted":  "Issue comment: deleted",
	"issue_comment_edited":   "Issue comment: edited",
	"issue_assigned":         "Issue: assigned",
	"issue_closed":           "Issue: closed",
	"issue_edited":           "Issue: edited",
	"issue_labeled":          "Issue: labeled",
	"issue_opened":           "Issue: opened",
	"issue_reopened":         "Issue: reopened",
	"issue_unassigned":       "Issue: unassigned",
	"issue_unlabeled":        "Issue: unlabeled",
	"pr_review_created":      "Pull request review comment",
	"pr_review_deleted":      "Pull request review comment: deleted",
	"pr_review_edited":       "Pull request review comment: edited",
	"pr_assigned":            "Pull request: assigned",
	"pr_closed":              "Pull request: closed",
	"pr_edited":              "Pull request: edited",
	"pr_labeled":             "Pull request: labeled",
	"pr_opened":              "Pull request: opened",
	"pr_reopened":            "Pull request: reopened",
	"pr_synchronize":         "Pull request: synchronize",
	"pr_unassigned":          "Pull request: unassigned",
	"pr_unlabeled":           "Pull request: unlabeled",
	"push":                   "Push",
	"release_published":      "Release published",
	"member_added":           "Repo: added collaborator",
	"team_add":               "Repo: added to a team",
	"fork":                   "Repo: forked",
	"public":                 "Repo: made public",
	"watch_started":          "Repo: starred",
	"gollum_created":         "Wiki: created page",
	"gollum_edited":          "Wiki: edited page",
}

func PayloadToString(event string, payloadOrigin interface{}) string {
	payload := payloadOrigin.(map[string]interface{})

	repository := payload["repository"].(map[string]interface{})
	repoName := repository["full_name"].(string)

	msg := fmt.Sprintf("[%s]", repoName)

	if event == "push" {
		pusher := payload["pusher"].(map[string]interface{})
		commits := payload["commits"].([]interface{})
		added := 0
		removed := 0
		modified := 0
		for _, commitOrigin := range commits {
			commit := commitOrigin.(map[string]interface{})
			added += len(commit["added"].([]interface{}))
			removed += len(commit["removed"].([]interface{}))
			modified += len(commit["modified"].([]interface{}))
		}
		msg = fmt.Sprintf("%s %s - pushed %d commit(s) to %s [+%d/-%d/\u00B1%d] \n %s", msg, pusher["name"], len(commits), strings.TrimLeft(payload["ref"].(string), "refs/heads/"), added, removed, modified, payload["compare"])
	} else if event == "issues" || event == "issue_comment" {
		sender := payload["sender"].(map[string]interface{})
		issue := payload["issue"].(map[string]interface{})
		msg = fmt.Sprintf("%s %s - %s action #%.0f: %s \n %s", msg, sender["login"], payload["action"], issue["number"], issue["title"], issue["html_url"])
	} else {
		sender := payload["sender"].(map[string]interface{})
		text := eventMsg[event]
		if text == "" {
			text = event
		}
		msg = fmt.Sprintf("%s %s - %s", msg, sender["login"], text)
	}
	return msg
}
