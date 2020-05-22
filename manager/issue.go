package manager

import (
	"github.com/google/go-github/v30/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/github-base/pkg/types"
	"github.com/pingcap/github-base/util"
)

// GetIssueByNumber get a issue by repo and number
func (mgr *Manager) GetIssueByNumber(repo *types.Repo, number int) (*types.Issue, error) {
	var issue types.Issue
	if err := mgr.storage.FindOne(&issue, "owner=? AND repo=? AND issue_number=?", repo.GetOwner(), repo.GetRepo(), number); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &issue, nil
}

// CreateIssue create from github issue
func (mgr *Manager) CreateIssue(repo *types.Repo, issue *github.Issue) error {
	mgr.Lock()
	defer mgr.Unlock()
	return mgr.CreateIssueNoLock(repo, issue)
}

// CreateIssueNoLock create from github issue without lock
func (mgr *Manager) CreateIssueNoLock(repo *types.Repo, issue *github.Issue) error {
	i, err := mgr.MakeIssuePatch(repo, issue)
	if err != nil {
		return errors.Trace(err)
	}
	if i == nil {
		return nil
	}
	return errors.Trace(mgr.UpdateIssue(i))
}

// UpdateIssue update from github issue
func (mgr *Manager) UpdateIssue(issue *types.Issue) error {
	return errors.Trace(mgr.storage.Save(issue))
}

// MakeIssuePatch make patch from issue
func (mgr *Manager) MakeIssuePatch(repo *types.Repo, issue *github.Issue) (*types.Issue, error) {
	i, err := mgr.GetIssueByNumber(repo, issue.GetNumber())
	if err == nil && i == nil {
		return mgr.MakeIssue(repo, issue)
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	if i.UpdatedAt.Equal(issue.GetUpdatedAt()) {
		return nil, nil
	}

	var labels []string
	for _, label := range issue.Labels {
		labels = append(labels, label.GetName())
	}
	status := issue.GetState()

	i.Title = issue.GetTitle()
	i.Body = issue.GetBody()
	i.Label = util.EncodeStringSlice(labels)
	i.Status = status
	i.UpdatedAt = issue.GetUpdatedAt()
	i.ClosedAt = issue.GetClosedAt()
	return i, nil
}

// MakeIssue make from github issue
func (mgr *Manager) MakeIssue(repo *types.Repo, issue *github.Issue) (*types.Issue, error) {
	isMember, err := mgr.isMember(issue.GetUser().GetLogin())
	if err != nil {
		return nil, errors.Trace(err)
	}
	relation := "member"
	if !isMember {
		relation = "not member"
	}
	var labels []string
	for _, label := range issue.Labels {
		labels = append(labels, label.GetName())
	}

	status := issue.GetState()

	i := types.Issue{
		Owner:       repo.GetOwner(),
		Repo:        repo.GetRepo(),
		Number:      issue.GetNumber(),
		Title:       issue.GetTitle(),
		Body:        issue.GetBody(),
		User:        issue.GetUser().GetLogin(),
		Association: "",
		Relation:    relation,
		Label:       util.EncodeStringSlice(labels),
		Status:      status,
		CreatedAt:   issue.GetCreatedAt(),
		UpdatedAt:   issue.GetUpdatedAt(),
		ClosedAt:    issue.GetClosedAt(),
	}
	return &i, nil
}
