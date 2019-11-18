package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/github-base/config"
	"github.com/pingcap/github-base/base/manager"
	globalManager "github.com/pingcap/github-base/manager"
	"github.com/pingcap/github-base/base/api"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	cfg          *config.Config
	configPath   string
)

func init() {
	flag.StringVar(&configPath, "c", "", "path to syncer config")
}

func main() {
	flag.Parse()

	cfg = config.GetGlobalConfig()
	if configPath != "" {
		err := cfg.Load(configPath)
		if err != nil {
			log.Fatalf(errors.ErrorStack(err))
		}
	}
	if err := cfg.Init(); err != nil {
		log.Error(err)
	}

	globalMgr, err := globalManager.New(cfg)
	if err != nil {
		log.Fatalf("can't run github-base: %v", errors.ErrorStack(err))
	}
	mgr := manager.New(globalMgr)

	go func() {
		log.Infof("begin to listen %s:%d 😄", cfg.Host, cfg.Port)
		app := iris.New()
		api.CreateRouter(app, mgr)
		if err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))); err != nil {
			log.Fatalf("app run error %v", err)
		}
	}()

	// syc, err := base.New(cfg, mgr)
	// if err != nil {
	// 	log.Fatalf("create syncer failed: %v", errors.ErrorStack(err))
	// }
	// syc.StartPolling()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sc
	log.Infof("Got signal %d to exit.", sig)
}
