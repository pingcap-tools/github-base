package api

import (
	// "github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/pingcap/github-base/base/manager"
)

// CreateRouter create router
func CreateRouter(app *iris.Application, mgr *manager.Manager) {
	hdl := newManagerHandler(mgr)
	// crs := cors.New(cors.Options{
	// 	AllowCredentials: true,
	// })
	// party := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	party := app.Party("/api")

	// ping
	party.Get("/ping", hdl.Ping)
	// github webhook
	party.Post("/webhook", hdl.Webhook)
}
