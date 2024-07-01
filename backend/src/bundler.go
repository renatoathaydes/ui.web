package src

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	esbuild "github.com/evanw/esbuild/pkg/api"
	set "github.com/hashicorp/go-set/v2"
)

type BuildContextResult struct {
	Name    string
	Err     *esbuild.ContextError
	Ctx     *esbuild.BuildContext
	IsError bool
}

type BundleOptions struct {
	BuildOpts *BuildOptions
	CommonDir string
}

const ModulesDir = "./modules"

// Bundle modules found in the given wrkdir.
//
// Only write out the files if write is true.
func Bundle(opts BundleOptions, write bool) ([]BuildContextResult, error) {
	modsDir := path.Join(opts.BuildOpts.Dir, ModulesDir)
	mods, err := CollectModules(modsDir)
	if err != nil {
		return nil, fmt.Errorf("problem collecting modules in directory %s: %v", modsDir, err)
	}
	return BundleModules(opts, mods, write)
}

// Bundle the given modules.
//
// Assumes that mods contains modules in `${opts.BuildOpts.Dir}/modules/`.
// Only write out the files if write is true.
func BundleModules(opts BundleOptions, mods []string, write bool) ([]BuildContextResult, error) {
	logger := opts.BuildOpts.Log()
	logger.Info("Bundling modules", "count", len(mods))
	modsDir := path.Join(opts.BuildOpts.Dir, ModulesDir)
	out := path.Join(modsDir, "out")
	_ = os.RemoveAll(out)
	mkDirErr := os.MkdirAll(out, os.ModePerm)
	if mkDirErr != nil {
		return nil, fmt.Errorf("unable to create frontend output folder %s due to %v", out, mkDirErr)
	}
	// esbuild requires an absolute path
	absWrkdir, err := filepath.Abs(opts.BuildOpts.Dir)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path modules in directory %s: %v", modsDir, err)
	}
	absModsDir := path.Join(absWrkdir, ModulesDir)
	commonMods, err := CollectModules(opts.CommonDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read common modules: %v", err)
	}

	plugins := makePlugins(mods, commonMods, absModsDir)

	channels := make(chan BuildContextResult, len(mods))
	var wg sync.WaitGroup
	wg.Add(len(mods))
	for i, mod := range mods {
		plugin := plugins[i]
		go func(mod string) {
			defer wg.Done()
			channels <- bundleModule(mod, absWrkdir, opts.BuildOpts.ForFrontend, plugin, write)
		}(mod)
	}
	wg.Wait()
	close(channels)
	var results []BuildContextResult
	for result := range channels {
		results = append(results, result)
	}
	return results, nil
}

func bundleModule(mod, wrkdir string, for_frontend bool, plugin []esbuild.Plugin, write bool) BuildContextResult {
	var platform esbuild.Platform
	var tree_shaking esbuild.TreeShaking
	if for_frontend {
		platform = esbuild.PlatformBrowser
		tree_shaking = esbuild.TreeShakingFalse
	} else {
		platform = esbuild.PlatformNode
		tree_shaking = esbuild.TreeShakingTrue
	}
	ctx, err := esbuild.Context(esbuild.BuildOptions{
		EntryPoints:   []string{path.Join(ModulesDir, mod)},
		Bundle:        true,
		Outdir:        path.Join(ModulesDir, "out"),
		Outbase:       ModulesDir,
		Write:         write,
		LogLevel:      esbuild.LogLevelWarning,
		AbsWorkingDir: wrkdir,
		TreeShaking:   tree_shaking,
		Format:        esbuild.FormatESModule,
		Platform:      platform,
		Plugins:       plugin,
	})
	return BuildContextResult{mod, err, &ctx, !reflect.ValueOf(err).IsNil()}
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

func makePlugins(mods, commonMods []string, wrkdir string) (res [][]esbuild.Plugin) {
	externalMods, modPaths := makeModuleSet(mods, commonMods)
	for i := range mods {
		plugin := makePlugin(modPaths[i], externalMods, wrkdir)
		res = append(res, []esbuild.Plugin{plugin})
	}
	return res
}

func makePlugin(mod string, externalMods *set.Set[string], wrkdir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "fix_imports",
		Setup: func(build esbuild.PluginBuild) {
			// modules at the root dir are marked external and their import paths are "fixed"
			// so that they don't include "modules/..."
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^(./|../)`},
				func(args esbuild.OnResolveArgs) (esbuild.OnResolveResult, error) {
					external := !strings.HasSuffix(args.Path, mod) &&
						externalMods.Contains(args.Path) &&
						isAtDir(args.Importer, wrkdir)
					slog.Debug("Checking import", "from", args.Importer, "path", args.Path, "isExternal", external)
					if external {
						return esbuild.OnResolveResult{
							Path:     ChangExtension(args.Path, ".js"),
							External: true,
						}, nil
					}
					return esbuild.OnResolveResult{
						External: false,
					}, nil
				})
		},
	}
}
