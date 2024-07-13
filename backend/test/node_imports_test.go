package test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"ui.web/server/src"
	"ui.web/server/src/logui"
)

func TestFindNodeModules(t *testing.T) {
	b := make([]byte, 0, 1024)
	w := bytes.NewBuffer(b)
	logger := slog.New(logui.New(slog.LevelDebug, os.Stdout, "TestFindNodeModules"))
	err := src.GenerateImportMaps(logger, "test_node_modules", w)
	require.Nil(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(w.Bytes(), &result)
	require.Nil(t, err)

	require.Equal(t, map[string]interface{}{
		"imports": map[string]interface{}{
			"codemirror":                  "./node_modules/codemirror/dist/index.js",
			"crelt":                       "./node_modules/crelt/index.js",
			"style-mod":                   "./node_modules/style-mod/src/style-mod.js",
			"typescript":                  "./node_modules/typescript/lib/typescript.js",
			"w3c-keyname":                 "./node_modules/w3c-keyname/index.js",
			"@codemirror/autocomplete":    "./node_modules/@codemirror/autocomplete/dist/index.js",
			"@codemirror/commands":        "./node_modules/@codemirror/commands/dist/index.js",
			"@codemirror/lang-javascript": "./node_modules/@codemirror/lang-javascript/dist/index.js",
			"@codemirror/language":        "./node_modules/@codemirror/language/dist/index.js",
			"@codemirror/lint":            "./node_modules/@codemirror/lint/dist/index.js",
			"@codemirror/search":          "./node_modules/@codemirror/search/dist/index.js",
			"@codemirror/state":           "./node_modules/@codemirror/state/dist/index.js",
			"@codemirror/view":            "./node_modules/@codemirror/view/dist/index.js",
			"@lezer/common":               "./node_modules/@lezer/common/dist/index.js",
			"@lezer/highlight":            "./node_modules/@lezer/highlight/dist/index.js",
			"@lezer/javascript":           "./node_modules/@lezer/javascript/dist/index.js",
			"@lezer/lr":                   "./node_modules/@lezer/lr/dist/index.js",
		},
	}, result)
}
