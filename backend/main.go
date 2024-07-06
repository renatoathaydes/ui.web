package main

import (
	"log"
	"log/slog"
	"os"
	"path"
	"sync"

	uiweb "ui.web/server/src"
	"ui.web/server/src/logui"
)

func main() {
	state := uiweb.State{}

	be := path.Join("..", "backend-js")
	fe := path.Join("..", "frontend")

	be_logger := slog.New(logui.New(slog.LevelDebug, os.Stdout, "be"))
	fe_logger := slog.New(logui.New(slog.LevelDebug, os.Stdout, "fe"))

	build(uiweb.BuildOptions{Dir: be, ForFrontend: false, Logger: be_logger}, &state)
	build(uiweb.BuildOptions{Dir: fe, ForFrontend: true, Logger: fe_logger}, &state)

	// run the process forever
	select {}
}

func build(build_opts uiweb.BuildOptions, state *uiweb.State) {
	modsDir := path.Join(build_opts.Dir, uiweb.ModulesDir)

	startup := func() chan bool {
		if build_opts.ForFrontend {
			return uiweb.StartServer(build_opts.Dir, build_opts.Logger, state)
		}
		bejs_logger := slog.New(logui.New(slog.LevelDebug, os.Stdout, "be-js"))
		bejs_modsDir := path.Join(build_opts.Dir, uiweb.ModulesDir)
		return uiweb.StartNode(bejs_modsDir, bejs_logger)
	}

	ok := uiweb.Build(build_opts)
	var restart struct {
		mutex *sync.Mutex
		ch    chan bool
	}
	restart.mutex = &sync.Mutex{}
	if ok {
		restart.ch = startup()
	} else {
		log.Fatalln("Could not build, check logs for errors!")
		return
	}

	uiweb.WatchAsync(modsDir, build_opts, func(build_opts uiweb.BuildOptions) bool {
		ok := uiweb.Build(build_opts)
		if ok {
			restart.mutex.Lock()
			defer restart.mutex.Unlock()
			select {
			case r := <-restart.ch:
				if r {
					restart.ch = startup()
				}
			default:
			}
		}
		return ok
	}, func() {
		build_opts.Log().Info("Stopping file watcher", "dir", modsDir)
	})
}
