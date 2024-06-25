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
	_, err := src.BundleModules("non-existent-modules", false)
	if err == nil {
		t.FailNow()
	}
	require.EqualError(t, err, "problem collecting modules in directory "+
		"non-existent-modules/modules: open non-existent-modules/modules: no such file or directory")
}

func TestBundler_EmptyDir(t *testing.T) {
	_, err := src.BundleModules("empty-modules", false)
	if err == nil {
		t.FailNow()
	}
	require.EqualError(t, err, "problem collecting modules in directory "+
		"empty-modules/modules: no frontend modules (*.mjs, *.mts) found")
}

func TestBundler_CollectOneTsModule(t *testing.T) {
	mods, err := src.CollectModules(path.Join("one-module", "modules"))
	require.Nil(t, err)
	require.Len(t, mods, 1)
	require.Equal(t, "hi.mts", mods[0])
}

func TestBundler_OneTsModule(t *testing.T) {
	t.Cleanup(func () {
		_ = os.RemoveAll("one-module/modules/out")
	})
	ctxs, err := src.BundleModules("one-module", false)
	require.Nil(t, err)
	require.Len(t, ctxs, 1)
	defer ctxs[0].Dispose()
	res := ctxs[0].Rebuild()
	assertOneOutput(t, res, "one-module/modules/out/hi.js")
}

func TestBundler_CollectManyTsModules(t *testing.T) {
	mods, err := src.CollectModules(path.Join("many-modules", "modules"))
	require.Nil(t, err)
	require.Len(t, mods, 3)
	require.Contains(t, mods, "one.mts", "two.mjs", "three.mts")
}

func TestBundler_ManyModules(t *testing.T) {
	t.Cleanup(func () {
		_ = os.RemoveAll("many-modules/modules/out")
	})
	ctxs, err := src.BundleModules("many-modules", false)
	require.Nil(t, err)
	require.Len(t, ctxs, 3)
	defer ctxs[0].Dispose()
	defer ctxs[1].Dispose()
	defer ctxs[2].Dispose()
	var results []esbuild.BuildResult
	for _, ctx := range ctxs {
		results = append(results, ctx.Rebuild())
	}
	require.Len(t, results[0].OutputFiles, 1)
	require.Len(t, results[1].OutputFiles, 1)
	require.Len(t, results[2].OutputFiles, 1)
	slices.SortFunc(results, compareByFirstOutput)
	assertOneOutput(t, results[0], "many-modules/modules/out/one.js")
	assertOneOutput(t, results[1], "many-modules/modules/out/three.js")
	assertOneOutput(t, results[2], "many-modules/modules/out/two.js")

	// this makes sure that we run esbuild correctly so that other files in modules/
	// are not bundled into modules that import it.
	assertModuleImportsModuleTwo(t, results[1])
}

func compareByFirstOutput(c1, c2 esbuild.BuildResult) int {
	return strings.Compare(c1.OutputFiles[0].Path, c2.OutputFiles[0].Path)
}

func assertOneOutput(t *testing.T, res esbuild.BuildResult, path string) {
	require.Len(t, res.Errors, 0)
	require.Len(t, res.OutputFiles, 1)
	expectedOutputFile, err2 := filepath.Abs(path)
	require.Nil(t, err2)
	require.Equal(t, expectedOutputFile, res.OutputFiles[0].Path)
}

func assertModuleImportsModuleTwo(t *testing.T, res esbuild.BuildResult) {
	js := string(res.OutputFiles[0].Contents)
	require.Contains(t, js, "from \"./two.mjs\"")
	require.NotContains(t, js, "Module two")
}
