package main

import (
	"log"
	"path"

	uiweb "ui.web/server/src"
)

func main() {
	state := uiweb.State{}

	fe := path.Join("..", "frontend")
	modsDir := path.Join(fe, uiweb.ModulesDir)

	// build in the background so we can start watching the dir
	go func() {
		ok := uiweb.Build(fe)
		if ok {
			uiweb.StartServer(fe, &state)
		}
	}()

	uiweb.WatchAsync(modsDir, fe, uiweb.Build, func() {
		log.Printf("Stopping file watcher for %s\n", modsDir)
	})

	// run the process forever
	select {}
}
