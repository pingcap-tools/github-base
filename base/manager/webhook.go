package manager

import (
	"github.com/google/go-github/github"
	"github.com/pingcap/github-base/pkg/types"
)

// Webhook process webhook
func (mgr *Manager)Webhook(repo *types.Repo, event interface{}) {
	switch event.(type) {
	case *github.PullRequestEvent:
		mgr.ProcessPullEvent(repo, event.(*github.PullRequestEvent))
	case *github.IssueCommentEvent:
		mgr.ProcessIssueCommentEvent(repo, event.(*github.IssueCommentEvent))
	case *github.PullRequestReviewCommentEvent:
		mgr.ProcessPullRequestReviewCommentEvent(repo, event.(*github.PullRequestReviewCommentEvent))
	case *github.PullRequestReviewEvent:
		mgr.ProcessPullRequestReviewEvent(repo, event.(*github.PullRequestReviewEvent))
	}
}
