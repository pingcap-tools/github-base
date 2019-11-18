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
	// _ = app.Party("/api", crs).AllowMethods(iris.MethodOptions)

	// ping
	app.Get("/ping", hdl.Ping)
	// github webhook
	app.Post("/webhook", hdl.Webhook)
}
