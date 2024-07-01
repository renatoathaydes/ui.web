package main

import (
	"log/slog"
	"os"
	"path"

	uiweb "ui.web/server/src"
)

func main() {
	state := uiweb.State{}

	be := path.Join("..", "backend-js")
	fe := path.Join("..", "frontend")

	be_logger := slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
		{Key: "from", Value: slog.StringValue("be")},
	}))
	fe_logger := slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
		{Key: "from", Value: slog.StringValue("fe")},
	}))

	build(uiweb.BuildOptions{Dir: be, ForFrontend: false, Logger: be_logger}, &state)
	build(uiweb.BuildOptions{Dir: fe, ForFrontend: true, Logger: fe_logger}, &state)

	// run the process forever
	select {}
}

func build(build_opts uiweb.BuildOptions, state *uiweb.State) {
	modsDir := path.Join(build_opts.Dir, uiweb.ModulesDir)

	ok := uiweb.Build(build_opts)
	if ok {
		if build_opts.ForFrontend {
			uiweb.StartServer(build_opts.Dir, state)
		} else {
			bejs_logger := slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
				{Key: "from", Value: slog.StringValue("be-js")},
			}))
			bejs_modsDir := path.Join(build_opts.Dir, uiweb.ModulesDir)
			uiweb.StartNode(bejs_modsDir, bejs_logger)
		}
	}

	uiweb.WatchAsync(modsDir, build_opts, uiweb.Build, func() {
		build_opts.Log().Info("Stopping file watcher", "dir", modsDir)
	})
}
