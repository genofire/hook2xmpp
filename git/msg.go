package git

import (
	"fmt"
	"strings"
)

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


type requestBody struct {
	Repository struct {
		HTMLURL  string `mapstructure:"html_url"`
		FullName string `mapstructure:"full_name"`
	} `mapstructure:"repository"`
	//push
	Pusher struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"pusher"`
	Commits []struct {
		Added    []interface{} `mapstructure:"added"`
		Removed  []interface{} `mapstructure:"removed"`
		Modified []interface{} `mapstructure:"modified"`
	} `mapstructure:"commits"`
	Compare string `mapstructure:"compare"`
	Ref     string `mapstructure:"ref"`
	// issue + fallback
	Sender struct {
		Login string `mapstructure:"login"`
	} `mapstructure:"sender"`
	// issue
	Action string `mapstructure:"action"`
	Issue  struct {
		HTMLURL string  `mapstructure:"html_url"`
		Number  float64 `mapstructure:"number"`
		Title   string  `mapstructure:"title"`
	} `mapstructure:"issue"`
}

func (r requestBody) String(event string) string {
	msg := fmt.Sprintf("[%s]", r.Repository.FullName)

	if event == "push" {
		added := 0
		removed := 0
		modified := 0
		for _, commit := range r.Commits {
			added += len(commit.Added)
			removed += len(commit.Removed)
			modified += len(commit.Modified)
		}
		msg = fmt.Sprintf("%s %s - pushed %d commit(s) to %s [+%d/-%d/\u00B1%d]: %s", msg, r.Pusher.Name, len(r.Commits), strings.TrimLeft(r.Ref, "refs/heads/"), added, removed, modified, r.Compare)
	} else if event == "issues" || event == "issue_comment" {
		msg = fmt.Sprintf("%s %s - %s action #%.0f: %s - %s", msg, r.Sender.Login, r.Action, r.Issue.Number, r.Issue.Title, r.Issue.HTMLURL)
	} else {
		text := eventMsg[event]
		if text == "" {
			text = event
		}
		msg = fmt.Sprintf("%s %s - %s", msg, r.Sender.Login, text)
	}
	return msg
}
