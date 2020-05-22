package manager

import (
	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/github-base/pkg/types"
)

// ProcessIssueEvent process issue event
func (mgr *Manager) ProcessIssueEvent(repo *types.Repo, event *github.IssueEvent) {
	mgr.ProcessIssue(repo, event.GetIssue())
}

// ProcessIssue process issue
func (mgr *Manager) ProcessIssue(repo *types.Repo, issue *github.Issue) {
	mgr.Lock()
	defer mgr.Unlock()
	issuePatch, err := mgr.mgr.MakeIssuePatch(repo, issue)
	if err != nil {
		log.Errorf("make issue patch failed %v", errors.ErrorStack(err))
		return
	}
	if issuePatch == nil {
		return
	}
	log.Infof("issue  %s/%s#%d", repo.GetOwner(), repo.GetRepo(), issuePatch.Number)
	if err := mgr.mgr.UpdateIssue(issuePatch); err != nil {
		log.Errorf("update issue failed %v", errors.ErrorStack(err))
	}
}
