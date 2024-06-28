package src

import (
	"log"
	"path"
)

func Build(fe string) bool {
	mods, err := CollectModules(path.Join(fe, ModulesDir))
	if err != nil {
		log.Printf("Bundler could not collect modules: %s\n", err.Error())
		return true
	}
	commonDir := path.Join(fe, "..", "common")
	ctxs, err := BundleModules(fe, commonDir, mods, true)
	if err != nil {
		log.Printf("Bundler failed: %s\n", err.Error())
		return true
	}
	for _, ctx := range ctxs {
		// TODO keep ctx between runs and re-use existing ones (which should be most)
		defer ctx.Dispose()
		res := ctx.Rebuild()
		if len(res.Errors) > 0 {
			log.Println("Bundler finished with errors, please check the logs above.")
		}
	}

	// TODO parse the template only once
	err = WriteEvalJs(mods, path.Join(fe, ModulesDir, "out", "eval.js"))
	if err != nil {
		log.Printf("Error writing eval.js: %s", err.Error())
	}

	err = CopyDir(path.Join(fe, "assets"),
		path.Join(fe, ModulesDir, "out"))
	if err != nil {
		log.Fatal("Could not copy assets", err)
	}

	return true
}
