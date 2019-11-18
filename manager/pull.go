package manager

import (
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/github-base/pkg/types"
	"github.com/pingcap/github-base/util"
)

// GetPullByNumber get a pull by repo and number
func (mgr *Manager)GetPullByNumber(repo *types.Repo, number int) (*types.Pull, error) {
	var pull types.Pull
	if err := mgr.storage.FindOne(&pull, "owner=? AND repo=? AND pull_number=?", repo.GetOwner(), repo.GetRepo(), number); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pull, nil
}

// CreatePull create from github pull
func (mgr *Manager)CreatePull(repo *types.Repo, pull *github.PullRequest) error {
	mgr.Lock()
	defer mgr.Unlock()
	return mgr.CreatePullNoLock(repo, pull)
}

// CreatePullNoLock create from github pull without lock
func (mgr *Manager)CreatePullNoLock(repo *types.Repo, pull *github.PullRequest) error {
	p, err := mgr.MakePullPatch(repo, pull)
	if err != nil {
		return errors.Trace(err)
	}
	if p == nil {
		return nil
	}
	return errors.Trace(mgr.UpdatePull(p))
}

// UpdatePull update from github pull
func (mgr *Manager)UpdatePull(pull *types.Pull) error {
	return errors.Trace(mgr.storage.Save(pull))
}

// MakePullPatch make patch from pull
func (mgr *Manager)MakePullPatch(repo *types.Repo, pull *github.PullRequest) (*types.Pull, error) {
	p, err := mgr.GetPullByNumber(repo, pull.GetNumber())
	if err == nil && p == nil {
		return mgr.MakePull(repo, pull)
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	if p.UpdatedAt.Equal(pull.GetUpdatedAt()) {
		return nil, nil
	}

	var labels []string
	for _, label := range pull.Labels {
		labels = append(labels, label.GetName())
	}
	status := pull.GetState()
	if pull.MergedAt != nil {
		status = "merged"
	}

	p.Title = pull.GetTitle()
	p.Body = pull.GetBody()
	p.Label = util.EncodeStringSlice(labels)
	p.Status = status
	p.UpdatedAt = pull.GetUpdatedAt()
	p.ClosedAt = pull.GetClosedAt()
	p.MergedAt = pull.GetMergedAt()
	return p, nil
}

// MakePull make from github pull
func (mgr *Manager)MakePull(repo *types.Repo, pull *github.PullRequest) (*types.Pull, error) {
	isMember, err := mgr.isMember(pull.GetUser().GetLogin())
	if err != nil {
		return nil, errors.Trace(err)
	}
	relation := "member"
	if !isMember {
		relation = "not member"
	}
	var labels []string
	for _, label := range pull.Labels {
		labels = append(labels, label.GetName())
	}

	status := pull.GetState()
	if pull.MergedAt != nil {
		status = "merged"
	}

	p := types.Pull{
		Owner: repo.GetOwner(),
		Repo: repo.GetRepo(),
		Number: pull.GetNumber(),
		Title: pull.GetTitle(),
		Body: pull.GetBody(),
		User: pull.GetUser().GetLogin(),
		Association: pull.GetAuthorAssociation(),
		Relation: relation,
		Label: util.EncodeStringSlice(labels),
		Status: status,
		CreatedAt: pull.GetCreatedAt(),
		UpdatedAt: pull.GetUpdatedAt(),
		ClosedAt: pull.GetClosedAt(),
		MergedAt: pull.GetMergedAt(),
	}
	return &p, nil
}
