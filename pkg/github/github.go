package github

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/ngaut/log"
	"github.com/pingcap/errors"
	"golang.org/x/oauth2"
	"os"
)

// GetGithubClient return client with auth
func GetGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		log.Errorf("get user info failed %v", errors.ErrorStack(err))
		os.Exit(1)
	}
	log.Infof("token user %s", user.GetLogin())
	return client
}
