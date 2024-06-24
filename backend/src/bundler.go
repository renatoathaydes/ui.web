package src

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

type buildContextResult struct {
	name string
	err  *esbuild.ContextError
	ctx  esbuild.BuildContext
}

const modulesDir = "./modules"
const modulesPattern = "./modules/*"

func BundleModules(wrkdir string, write bool) ([]esbuild.BuildContext, error) {
	modsDir := path.Join(wrkdir, modulesDir)
	mods, err := CollectModules(modsDir)
	if err != nil {
		return nil, fmt.Errorf("problem collecting modules in directory %s: %v", modsDir, err)
	}
	out := path.Join(wrkdir, "./out")
	_ = os.RemoveAll(out)
	mkDirErr := os.MkdirAll(out, os.ModePerm)
	if mkDirErr != nil {
		return nil, fmt.Errorf("unable to create frontend output folder %s due to %v", out, mkDirErr)
	}
	// esbuild requires an absolute path
	absWrkdir, err := filepath.Abs(wrkdir)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path modules in directory %s: %v", modsDir, err)
	}

	results := make(chan buildContextResult, len(mods))
	var wg sync.WaitGroup
	wg.Add(len(mods))
	for _, mod := range mods {
		go func(mod string) {
			defer wg.Done()
			results <- bundleModule(mod, absWrkdir, out, write)
		}(mod)
	}
	wg.Wait()
	close(results)
	var success []esbuild.BuildContext
	var errors []buildContextResult
	for result := range results {
		if !reflect.ValueOf(result.err).IsNil() {
			errors = append(errors, result)
		} else {
			success = append(success, result.ctx)
		}
	}
	finalError := errorIn(errors)
	if finalError != nil {
		return nil, finalError
	}
	return success, nil
}

func errorIn(results []buildContextResult) error {
	if len(results) == 0 {
		return nil
	}
	message := "Failed to bundle modules:\n"
	for _, e := range results {
		message += fmt.Sprintf("  * %s - %s\n", e.name, e.err.Error())
	}
	return errors.New(message)
}

func bundleModule(mod, wrkdir, out string, write bool) buildContextResult {
	ctx, err := esbuild.Context(esbuild.BuildOptions{
		EntryPoints:   []string{mod},
		Bundle:        true,
		Outdir:        path.Base(out),
		Write:         write,
		LogLevel:      esbuild.LogLevelError,
		AbsWorkingDir: wrkdir,
		TreeShaking:   esbuild.TreeShakingFalse,
		Format:        esbuild.FormatESModule,
		// all modules are included by the HTML and hence should not be embedded
		External: []string{modulesPattern},
	})
	return buildContextResult{mod, err, ctx}
}

func CollectModules(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range files {
		name := entry.Name()
		if entry.Type().IsRegular() &&
			(strings.HasSuffix(name, ".mts") ||
				strings.HasSuffix(name, ".mjs")) {
			result = append(result, path.Join(modulesDir, name))
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no frontend modules (*.mjs, *.mts) found")
	}
	return result, nil
}
