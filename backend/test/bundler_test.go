package test

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"ui.web/server/src"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func TestBundler_NonExistentDir(t *testing.T) {
	build_opts := src.BuildOptions{Dir: path.Join("assets", "non-existent-modules")}
	opts := src.BundleOptions{BuildOpts: &build_opts, CommonDir: path.Join("assets", "common")}
	_, err := src.Bundle(opts, false)
	if err == nil {
		t.FailNow()
	}
	dir := path.Join("assets", "non-existent-modules", "modules")
	require.EqualError(t, err, "problem collecting modules in directory "+
		dir+": open "+dir+": no such file or directory")
}

func TestBundler_EmptyDir(t *testing.T) {
	build_opts := src.BuildOptions{Dir: path.Join("assets", "empty-modules")}
	opts := src.BundleOptions{BuildOpts: &build_opts, CommonDir: path.Join("assets", "common")}
	_, err := src.Bundle(opts, false)
	if err == nil {
		t.FailNow()
	}
	dir := path.Join("assets", "empty-modules", "modules")
	require.EqualError(t, err, "problem collecting modules in directory "+
		dir+": no frontend modules (*.mjs, *.mts) found")
}

func TestBundler_CollectOneTsModule(t *testing.T) {
	mods, err := src.CollectModules(path.Join("assets", "one-module", "modules"))
	require.Nil(t, err)
	require.Len(t, mods, 1)
	require.Equal(t, "hi.mts", mods[0])
}

func TestBundler_OneTsModule(t *testing.T) {
	module := path.Join("assets", "one-module")
	t.Cleanup(func() {
		_ = os.RemoveAll(path.Join(module, "modules", "out"))
	})
	build_opts := src.BuildOptions{Dir: module, ForFrontend: true}
	opts := src.BundleOptions{BuildOpts: &build_opts, CommonDir: path.Join("assets", "common")}
	results, err := src.Bundle(opts, false)
	require.Nil(t, err)
	require.Len(t, results, 1)
	require.Nil(t, results[0].Err)
	ctx := *results[0].Ctx
	defer ctx.Dispose()
	res := ctx.Rebuild()
	assertOneOutput(t, res, path.Join(module, "modules", "out", "hi.js"))
}

func TestBundler_CollectManyTsModules(t *testing.T) {
	mods, err := src.CollectModules(path.Join("assets", "many-modules", "modules"))
	require.Nil(t, err)
	require.Len(t, mods, 3)
	require.Contains(t, mods, "one.mjs", "two.mts", "three.mts")
}

func TestBundler_ManyModules(t *testing.T) {
	module := path.Join("assets", "many-modules")
	t.Cleanup(func() {
		_ = os.RemoveAll(path.Join(module, "modules", "out"))
	})
	build_opts := src.BuildOptions{Dir: module, ForFrontend: true}
	opts := src.BundleOptions{BuildOpts: &build_opts, CommonDir: path.Join("assets", "common")}
	results, err := src.Bundle(opts, false)
	require.Nil(t, err)
	require.Len(t, results, 3)
	ctxs := make([]esbuild.BuildContext, 3)
	for i := 0; i < 3; i++ {
		ctxs[i] = *results[i].Ctx
	}
	defer ctxs[0].Dispose()
	defer ctxs[1].Dispose()
	defer ctxs[2].Dispose()
	var build_results []esbuild.BuildResult
	for _, ctx := range ctxs {
		build_results = append(build_results, ctx.Rebuild())
	}
	require.Len(t, build_results[0].OutputFiles, 1)
	require.Len(t, build_results[1].OutputFiles, 1)
	require.Len(t, build_results[2].OutputFiles, 1)

	slices.SortFunc(build_results, compareByFirstOutput)
	assertOneOutput(t, build_results[0], path.Join(module, "modules", "out", "one.js"))
	assertOneOutput(t, build_results[1], path.Join(module, "modules", "out", "three.js"))
	assertOneOutput(t, build_results[2], path.Join(module, "modules", "out", "two.js"))

	// this makes sure that we run esbuild correctly so that other files in modules/
	// are not bundled into modules that import it.
	assertModuleImportsModuleTwo(t, build_results[1])
	assertModuleImportsCommonModule(t, build_results[1])
}

func compareByFirstOutput(c1, c2 esbuild.BuildResult) int {
	return strings.Compare(c1.OutputFiles[0].Path, c2.OutputFiles[0].Path)
}

func assertOneOutput(t *testing.T, res esbuild.BuildResult, path string) {
	require.Len(t, res.Errors, 0)
	require.Len(t, res.OutputFiles, 1)
	expectedOutputFile, err := filepath.Abs(path)
	require.Nil(t, err)
	require.Equal(t, expectedOutputFile, res.OutputFiles[0].Path)
}

func assertModuleImportsModuleTwo(t *testing.T, res esbuild.BuildResult) {
	js := string(res.OutputFiles[0].Contents)
	require.Contains(t, js, "from \"./two.js\"")
	require.NotContains(t, js, "Module two")
}

func assertModuleImportsCommonModule(t *testing.T, res esbuild.BuildResult) {
	js := string(res.OutputFiles[0].Contents)
	require.Contains(t, js, "from \"../../common/common.js\"")
	require.NotContains(t, js, "common module code")
}
