package manager

import (
	"github.com/google/go-github/v30/github"
	"github.com/pingcap/github-base/pkg/types"
)

// Webhook process webhook
func (mgr *Manager) Webhook(repo *types.Repo, event interface{}) {
	switch event := event.(type) {
	case *github.PullRequestEvent:
		mgr.ProcessPullEvent(repo, event)
	case *github.IssueEvent:
		mgr.ProcessIssueEvent(repo, event)
	case *github.IssueCommentEvent:
		mgr.ProcessIssueCommentEvent(repo, event)
	case *github.PullRequestReviewCommentEvent:
		mgr.ProcessPullRequestReviewCommentEvent(repo, event)
	case *github.PullRequestReviewEvent:
		mgr.ProcessPullRequestReviewEvent(repo, event)
	}
}
