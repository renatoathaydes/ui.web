package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
	uiweb "ui.web/server/src"
)

func main() {
	wrkdir, err := filepath.Abs("../frontend")
	if err != nil {
		log.Fatal("Frontend.AbsolutePath", err)
	}
	feDist := path.Join(wrkdir, "./frontend-dist")
	_ = os.RemoveAll(feDist)
	err = os.MkdirAll(feDist, os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create frontend folder", err)
	}

	ctx, ctxErr := esbuild.Context(esbuild.BuildOptions{
		EntryPoints:   []string{"index.js"},
		Bundle:        true,
		Outdir:        path.Base(feDist),
		Write:         true,
		LogLevel:      esbuild.LogLevelError,
		AbsWorkingDir: wrkdir,
	})
	if ctxErr != nil {
		log.Fatal("esbuild.Context", ctxErr)
	}
	err = ctx.Watch(esbuild.WatchOptions{})
	if err != nil {
		log.Fatal("esbuild.Watch", err)
	}

	uiweb.CopyFile(path.Join(wrkdir, "assets", "index.html"), path.Join(feDist, "index.html"))
	uiweb.StartServer(feDist)

	// run the process forever
	select {}
}
