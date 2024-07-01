package src

import (
	"log/slog"
	"path"
)

type BuildOptions struct {
	Dir         string
	ForFrontend bool
	Logger      *slog.Logger
}

func (b BuildOptions) Log() *slog.Logger {
	l := b.Logger
	if l != nil {
		return l
	}
	return slog.Default()
}

func Build(options BuildOptions) bool {
	logger := options.Log()
	dir := options.Dir
	mods, err := CollectModules(path.Join(dir, ModulesDir))
	if err != nil {
		logger.Warn("Bundler could not collect modules", "error", err)
		return true
	}
	opts := BundleOptions{BuildOpts: &options, CommonDir: path.Join(dir, "..", "common")}
	results, err := BundleModules(opts, mods, true)
	if err != nil {
		logger.Error("Bundler failed", "error", err)
		return true
	}
	all_success := true
	for _, res := range results {
		// TODO keep ctx between runs and re-use existing ones (which should be most)
		if res.IsError {
			logger.Error("There was an error creating the context", "module", res.Name, "error", res.Err)
			return false
		}
		ctx := *res.Ctx
		defer ctx.Dispose()
		build_res := ctx.Rebuild()
		if len(build_res.Errors) > 0 {
			logger.Error("Module was not bundled due to errors", "module", res.Name)
			all_success = false
		} else {
			logger.Info("Module bundled successfully", "module", res.Name)
		}
	}
	if all_success {
		logger.Info("All modules bundled successfully.")
	} else {
		logger.Warn("Bundler finished with errors, please check the logs above.")
	}

	// TODO parse the template only once
	err = WriteEvalJs(mods, path.Join(dir, ModulesDir, "out", "eval.js"), options.ForFrontend)
	if err != nil {
		logger.Warn("Error writing eval.js", "error", err)
	}

	err = CopyDir(path.Join(dir, "assets"),
		path.Join(dir, ModulesDir, "out"))
	if err != nil {
		logger.Warn("Could not copy assets", "error", err)
	}

	return true
}
