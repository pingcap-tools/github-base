package manager

import (
	"github.com/ngaut/log"
	"github.com/google/go-github/github"
	"github.com/pingcap/github-base/pkg/types"
)

// ProcessIssueCommentEvent process issue comment event
func (mgr *Manager)ProcessIssueCommentEvent(repo *types.Repo, event *github.IssueCommentEvent) {
	mgr.ProcessIssueComment(repo, event.GetIssue(), event.GetComment())	
}

// ProcessPullRequestReviewCommentEvent process pull request review comment event
func (mgr *Manager)ProcessPullRequestReviewCommentEvent(repo *types.Repo, event *github.PullRequestReviewCommentEvent) {
	mgr.ProcessPullRequestReviewComment(repo, event.GetPullRequest(), event.GetComment())
}

// ProcessPullRequestReviewEvent process pull request review event
func (mgr *Manager)ProcessPullRequestReviewEvent(repo *types.Repo, event *github.PullRequestReviewEvent) {
	mgr.ProcessPullRequestReview(repo, event.GetPullRequest(), event.GetReview())
}

// ProcessIssueComment process issue comment
func (mgr *Manager)ProcessIssueComment(repo *types.Repo, issue *github.Issue, comment *github.IssueComment) {
	patch, err := mgr.mgr.MakeCommentPatch(repo, comment, &types.CommentAttach{
		CommentType: "common comment",
		Number: issue.GetNumber(),
		CreatedAt: comment.GetCreatedAt(),
		UpdatedAt: comment.GetUpdatedAt(),
		Association: comment.GetAuthorAssociation(),
	})
	if err != nil {
		log.Errorf("create patch failed %v", err)
		return
	} else if patch == nil {
		return
	}
	if err := mgr.mgr.UpdateComment(patch); err != nil {
		log.Errorf("update patch failed %v", err)
	}
}

// ProcessPullRequestReviewComment process pull request review comment
func (mgr *Manager)ProcessPullRequestReviewComment(repo *types.Repo, pull *github.PullRequest, comment *github.PullRequestComment) {
	patch, err := mgr.mgr.MakeCommentPatch(repo, comment, &types.CommentAttach{
		CommentType: "review comment",
		Number: pull.GetNumber(),
		CreatedAt: comment.GetCreatedAt(),
		UpdatedAt: comment.GetUpdatedAt(),
		Association: comment.GetAuthorAssociation(),
	})
	if err != nil {
		log.Errorf("create patch failed %v", err)
		return
	} else if patch == nil {
		return
	}
	if err := mgr.mgr.UpdateComment(patch); err != nil {
		log.Errorf("update patch failed %v", err)
	}	
}

// ProcessPullRequestReview process pull request review
func (mgr *Manager)ProcessPullRequestReview(repo *types.Repo, pull *github.PullRequest, review *github.PullRequestReview) {
	patch, err := mgr.mgr.MakeCommentPatch(repo, pull, &types.CommentAttach{
		CommentType: "review",
		Number: pull.GetNumber(),
		CreatedAt: review.GetSubmittedAt(),
		UpdatedAt: review.GetSubmittedAt(),
		Association: "",
	})
	if err != nil {
		log.Errorf("create patch failed %v", err)
		return
	} else if patch == nil {
		return
	}
	if err := mgr.mgr.UpdateComment(patch); err != nil {
		log.Errorf("update patch failed %v", err)
	}		
}