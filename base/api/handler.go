package api

import (
	"github.com/kataras/iris"
	"github.com/pingcap/github-base/base/manager"
)

// ManagerHandler is manager api handler
type ManagerHandler struct {
	mgr *manager.Manager
}

func newManagerHandler(mgr *manager.Manager) *ManagerHandler {
	return &ManagerHandler{
		mgr: mgr,
	}
}

// Ping tests service available
func (hdl *ManagerHandler)Ping(ctx iris.Context) {
	ctx.WriteString("pong")
}
