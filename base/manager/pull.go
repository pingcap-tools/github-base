package manager

import (
	"github.com/ngaut/log"
	"github.com/juju/errors"
	"github.com/google/go-github/github"
	"github.com/pingcap/github-base/pkg/types"
)

// ProcessPullEvent process pull event
func (mgr *Manager) ProcessPullEvent(repo *types.Repo, event *github.PullRequestEvent) {
	mgr.ProcessPull(repo, event.GetPullRequest())
}

// ProcessPull process pull
func (mgr *Manager) ProcessPull(repo *types.Repo, pull *github.PullRequest) {
	pullPatch, err := mgr.mgr.MakePullPatch(repo, pull)
	if err != nil {
		log.Errorf("make pull patch failed %v", errors.ErrorStack(err))
		return
	}
	if pullPatch == nil {
		return
	}
	log.Infof("pull  %s/%s#%d", repo.GetOwner(), repo.GetRepo(), pullPatch.GetNumber())
	if err := mgr.mgr.UpdatePull(pullPatch); err != nil {
		log.Errorf("update pull failed %v", errors.ErrorStack(err))
	}
}
