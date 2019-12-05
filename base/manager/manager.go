package manager

import (
	"sync"
	"github.com/pingcap/github-base/config"
	"github.com/pingcap/github-base/manager"
)

// Manager struct
type Manager struct {
	sync.Mutex
	mgr *manager.Manager
}

// New create manager
func New(mgr *manager.Manager) *Manager {
	return &Manager{
		mgr: mgr,
	}
}

// GetConfig return config struct
func (mgr *Manager)GetConfig() *config.Config {
	return mgr.mgr.Config
}
