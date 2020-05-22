package manager

import (
	"context"
	"time"

	"github.com/google/go-github/v30/github"
	"github.com/ngaut/log"
	"github.com/pingcap/errors"
	"github.com/pingcap/github-base/pkg/types"
	"github.com/pingcap/github-base/util"
)

// PollingOwner poll pulls from user or org
func (mgr *Manager) PollingOwner(owner string) {
	var (
		page    = 0
		perpage = 100
		batch   []*github.Repository
		res     *github.Response
		err     error
	)

	for page == 0 || len(batch) == perpage {
		page++
		if err := util.RetryOnError(context.Background(), 3, func() error {
			ctx, cancel := util.TimeoutContext(time.Minute)
			defer cancel()

			batch, res, err = mgr.mgr.Github.Repositories.ListByOrg(ctx, owner, &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perpage,
				},
			})
			if err != nil {
				return errors.Trace(err)
			}
			for _, repo := range batch {
				mgr.PollingRepo(&types.Repo{
					Owner: owner,
					Repo:  repo.GetName(),
				})
			}
			return nil
		}); err != nil {
			log.Errorf("fetch repos error %v", errors.ErrorStack(err))
		}
		if res.Rate.Limit < 500 {
			time.Sleep(res.Rate.Reset.Sub(time.Now()) + 5*time.Minute)
		}
	}
}

// PollingRepo poll pulls from repo
func (mgr *Manager) PollingRepo(repo *types.Repo) {
	var (
		page    = 0
		perpage = 100
		batch   []*github.PullRequest
		res     *github.Response
		err     error
	)

	for page == 0 || len(batch) == perpage {
		page++
		if err := util.RetryOnError(context.Background(), 3, func() error {
			ctx, cancel := util.TimeoutContext(time.Minute)
			defer cancel()

			batch, res, err = mgr.mgr.Github.PullRequests.List(ctx, repo.Owner, repo.Repo, &github.PullRequestListOptions{
				State:     "all",
				Direction: "asc",
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perpage,
				},
			})
			time.Sleep(time.Second)
			if err != nil {
				return errors.Trace(err)
			}
			for _, pull := range batch {
				mgr.ProcessPull(repo, pull)
				mgr.PollingPullRequestReviewComment(repo, pull)
				mgr.PollingPullRequestReview(repo, pull)
				mgr.PollingIssueComment(repo, pull2issue(pull))
			}
			return nil
		}); err != nil {
			log.Errorf("fetch pulls error %v", errors.ErrorStack(err))
		}
		if res.Rate.Limit < 500 {
			time.Sleep(res.Rate.Reset.Sub(time.Now()) + 5*time.Minute)
		}
	}
}

// PollingIssueComment poll comments of an issue
func (mgr *Manager) PollingIssueComment(repo *types.Repo, issue *github.Issue) {
	var (
		page    = 0
		perpage = 100
		batch   []*github.IssueComment
		res     *github.Response
		err     error
	)

	for page == 0 || len(batch) == perpage {
		page++
		if err := util.RetryOnError(context.Background(), 3, func() error {
			ctx, cancel := util.TimeoutContext(time.Minute)
			defer cancel()

			batch, res, err = mgr.mgr.Github.Issues.ListComments(ctx, repo.Owner, repo.Repo, issue.GetNumber(), &github.IssueListCommentsOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perpage,
				},
			})
			time.Sleep(time.Second)
			if err != nil {
				return errors.Trace(err)
			}
			for _, comment := range batch {
				mgr.ProcessIssueComment(repo, issue, comment)
			}
			return nil
		}); err != nil {
			log.Errorf("fetch issue comments error %v", errors.ErrorStack(err))
		}
		if res.Rate.Limit < 500 {
			time.Sleep(res.Rate.Reset.Sub(time.Now()) + 5*time.Minute)
		}
	}
}

// PollingPullRequestReviewComment poll all review comments of a pull request
func (mgr *Manager) PollingPullRequestReviewComment(repo *types.Repo, pull *github.PullRequest) {
	var (
		page    = 0
		perpage = 100
		batch   []*github.PullRequestComment
		res     *github.Response
		err     error
	)

	for page == 0 || len(batch) == perpage {
		page++
		if err := util.RetryOnError(context.Background(), 3, func() error {
			ctx, cancel := util.TimeoutContext(time.Minute)
			defer cancel()

			batch, res, err = mgr.mgr.Github.PullRequests.ListComments(ctx, repo.Owner, repo.Repo, pull.GetNumber(), &github.PullRequestListCommentsOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perpage,
				},
			})
			time.Sleep(time.Second)
			if err != nil {
				return errors.Trace(err)
			}
			for _, comment := range batch {
				mgr.ProcessPullRequestReviewComment(repo, pull, comment)
			}
			return nil
		}); err != nil {
			log.Errorf("fetch review comments error %v", errors.ErrorStack(err))
		}
		if res.Rate.Limit < 500 {
			time.Sleep(res.Rate.Reset.Sub(time.Now()) + 5*time.Minute)
		}
	}
}

// PollingPullRequestReview poll all reviews of a pull request
func (mgr *Manager) PollingPullRequestReview(repo *types.Repo, pull *github.PullRequest) {
	var (
		page    = 0
		perpage = 100
		batch   []*github.PullRequestReview
		res     *github.Response
		err     error
	)

	for page == 0 || len(batch) == perpage {
		page++
		if err := util.RetryOnError(context.Background(), 3, func() error {
			ctx, cancel := util.TimeoutContext(time.Minute)
			defer cancel()

			batch, res, err = mgr.mgr.Github.PullRequests.ListReviews(ctx, repo.Owner, repo.Repo, pull.GetNumber(), &github.ListOptions{
				Page:    page,
				PerPage: perpage,
			})
			time.Sleep(time.Second)
			if err != nil {
				return errors.Trace(err)
			}
			for _, comment := range batch {
				mgr.ProcessPullRequestReview(repo, pull, comment)
			}
			return nil
		}); err != nil {
			log.Errorf("fetch reviews error %v", errors.ErrorStack(err))
		}
		if res.Rate.Limit < 500 {
			time.Sleep(res.Rate.Reset.Sub(time.Now()) + 5*time.Minute)
		}
	}
}

// issue2pull transfer issue to pull with some common fields
// for those fields exist in pulls only
// you should get request by another API call
func issue2pull(issue *github.Issue) *github.PullRequest {
	if !issue.IsPullRequest() {
		return nil
	}

	pull := github.PullRequest{
		ID:          issue.ID,
		Number:      issue.Number,
		State:       issue.State,
		Locked:      issue.Locked,
		Title:       issue.Title,
		Body:        issue.Body,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
		ClosedAt:    issue.ClosedAt,
		Labels:      issue.Labels,
		User:        issue.User,
		URL:         issue.URL,
		HTMLURL:     issue.HTMLURL,
		IssueURL:    issue.URL,
		CommentsURL: issue.CommentsURL,
		Milestone:   issue.Milestone,
		NodeID:      issue.NodeID,
	}
	return &pull
}

// pull2issue transfer pull to issue with some common fields
// for those fields exist in issues only
// you should get request by another API call
func pull2issue(pull *github.PullRequest) *github.Issue {

	issue := github.Issue{
		ID:          pull.ID,
		Number:      pull.Number,
		State:       pull.State,
		Locked:      pull.Locked,
		Title:       pull.Title,
		Body:        pull.Body,
		CreatedAt:   pull.CreatedAt,
		UpdatedAt:   pull.UpdatedAt,
		ClosedAt:    pull.ClosedAt,
		Labels:      pull.Labels,
		User:        pull.User,
		URL:         pull.URL,
		HTMLURL:     pull.HTMLURL,
		CommentsURL: pull.CommentsURL,
		Milestone:   pull.Milestone,
		NodeID:      pull.NodeID,
	}

	return &issue
}
