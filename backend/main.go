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
	state := uiweb.State{}
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
		EntryPoints:   []string{"entrypoint.mts", "eval.mts"},
		Bundle:        true,
		Outdir:        path.Base(feDist),
		Write:         true,
		LogLevel:      esbuild.LogLevelError,
		AbsWorkingDir: wrkdir,
		TreeShaking:   esbuild.TreeShakingFalse,
		Format:        esbuild.FormatESModule,
		// all modules are included by the HTML and hence should not be embedded
		External: []string{"./modules/*"},
	})
	if ctxErr != nil {
		log.Fatal("esbuild.Context", ctxErr)
	}
	err = ctx.Watch(esbuild.WatchOptions{})
	if err != nil {
		log.Fatal("esbuild.Watch", err)
	}

	uiweb.CopyFile(path.Join(wrkdir, "assets", "index.html"), path.Join(feDist, "index.html"))
	uiweb.StartServer(feDist, &state)

	// run the process forever
	select {}
}
