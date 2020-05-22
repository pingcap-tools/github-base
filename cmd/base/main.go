package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
	"github.com/pingcap/github-base/base/api"
	"github.com/pingcap/github-base/base/manager"
	"github.com/pingcap/github-base/config"
	globalManager "github.com/pingcap/github-base/manager"
	"github.com/pingcap/github-base/pkg/types"
)

var (
	cfg        *config.Config
	configPath string
	repo       string
	owner      string
)

func init() {
	flag.StringVar(&configPath, "c", "", "path to syncer config")
	flag.StringVar(&repo, "r", "", "polling repo")
	flag.StringVar(&owner, "o", "", "polling org")
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

	if repo != "" {
		if owner == "" {
			log.Fatal("org should not be empty")
		}
		mgr.PollingRepo(&types.Repo{Owner: owner, Repo: repo})
		return
	}
	if owner != "" {
		mgr.PollingOwner(owner)
		return
	}

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
