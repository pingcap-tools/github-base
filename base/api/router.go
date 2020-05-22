package api

import (
	// "github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/pingcap/github-base/base/manager"
)

// CreateRouter create router
func CreateRouter(app *iris.Application, mgr *manager.Manager) {
	hdl := newManagerHandler(mgr)

	party := app.Party("/api")

	// ping
	party.Get("/ping", hdl.Ping)
	// github webhook
	party.Post("/webhook", hdl.Webhook)
}
