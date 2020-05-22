package manager

import (
	"sync"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/pingcap/github-base/config"
	githubInit "github.com/pingcap/github-base/pkg/github"
	"github.com/pingcap/github-base/pkg/storage"
	"github.com/pingcap/github-base/pkg/storage/basic"
)

// Manager represent schrodinger syncer
type Manager struct {
	sync.Mutex
	Config  *config.Config
	storage basic.Driver
	Github  *github.Client
	Members map[string]bool
}

// New init manager
func New(cfg *config.Config) (*Manager, error) {
	s, err := storage.New(cfg.Database)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Manager{
		Config:  cfg,
		storage: s,
		Github:  githubInit.GetGithubClient(cfg.GithubToken),
		Members: make(map[string]bool),
	}, nil
}
