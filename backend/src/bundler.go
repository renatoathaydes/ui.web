package src

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	esbuild "github.com/evanw/esbuild/pkg/api"
	set "github.com/hashicorp/go-set/v2"
)

type buildContextResult struct {
	name string
	err  *esbuild.ContextError
	ctx  esbuild.BuildContext
}

const ModulesDir = "./modules"

// Bundle modules found in the given wrkdir.
//
// Only write out the files if write is true.
func Bundle(wrkdir, commonDir string, write bool) ([]esbuild.BuildContext, error) {
	modsDir := path.Join(wrkdir, ModulesDir)
	mods, err := CollectModules(modsDir)
	if err != nil {
		return nil, fmt.Errorf("problem collecting modules in directory %s: %v", modsDir, err)
	}
	return BundleModules(wrkdir, commonDir, mods, write)
}

// Bundle the given modules.
//
// Assumes that mods contains modules in `wrkdir/modules/`.
// Only write out the files if write is true.
func BundleModules(wrkdir, commonDir string, mods []string, write bool) ([]esbuild.BuildContext, error) {
	log.Printf("Building %d module(s).\n", len(mods))
	modsDir := path.Join(wrkdir, ModulesDir)
	out := path.Join(modsDir, "out")
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

	commonMods, err := CollectModules(commonDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read common modules: %v", err)
	}

	plugins := makePlugins(mods, commonMods)

	results := make(chan buildContextResult, len(mods))
	var wg sync.WaitGroup
	wg.Add(len(mods))
	for i, mod := range mods {
		plugin := plugins[i]
		go func(mod string) {
			defer wg.Done()
			results <- bundleModule(mod, absWrkdir, plugin, write)
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

func bundleModule(mod, wrkdir string, plugin esbuild.Plugin, write bool) buildContextResult {
	ctx, err := esbuild.Context(esbuild.BuildOptions{
		EntryPoints:   []string{path.Join(ModulesDir, mod)},
		Bundle:        true,
		Outdir:        path.Join(ModulesDir, "out"),
		Outbase:       ModulesDir,
		Write:         write,
		LogLevel:      esbuild.LogLevelError,
		AbsWorkingDir: wrkdir,
		TreeShaking:   esbuild.TreeShakingFalse,
		Format:        esbuild.FormatESModule,
		Plugins:       []esbuild.Plugin{plugin},
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
			result = append(result, name)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no frontend modules (*.mjs, *.mts) found")
	}
	return result, nil
}

func makeModuleSet(mods, commonMods []string) (externalMods *set.Set[string], modPaths []string) {
	externalMods = set.New[string](len(mods) + len(commonMods))
	modPaths = make([]string, len(mods))
	for i, mod := range mods {
		externalMods.Insert(fmt.Sprintf("./%s", mod))
		modPaths[i] = fmt.Sprintf("%s/%s", ModulesDir, mod)
	}
	for _, mod := range commonMods {
		externalMods.Insert(fmt.Sprintf("../../common/%s", mod))
	}
	return
}

func makePlugins(mods, commonMods []string) []esbuild.Plugin {
	externalMods, modPaths := makeModuleSet(mods, commonMods)
	var res []esbuild.Plugin
	for i := range mods {
		res = append(res, makePlugin(modPaths[i], externalMods))
	}
	return res
}

func makePlugin(mod string, externalMods *set.Set[string]) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "fix_imports",
		Setup: func(build esbuild.PluginBuild) {
			// modules at the root dir are marked external and their import paths are "fixed"
			// so that they don't include "modules/..."
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^(./|../)`},
				func(args esbuild.OnResolveArgs) (esbuild.OnResolveResult, error) {
					external := !strings.HasSuffix(args.Path, mod) && externalMods.Contains(args.Path)
					if external {
						return esbuild.OnResolveResult{
							Path:     ChangExtension(args.Path, ".js"),
							External: true,
						}, nil
					}
					return esbuild.OnResolveResult{
						Path:     filepath.Join(args.ResolveDir, args.Path),
						External: false,
					}, nil

				})
		},
	}
}
