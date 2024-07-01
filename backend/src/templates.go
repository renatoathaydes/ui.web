package src

import (
	"io"
	"os"
	"path"
	"text/template"
)

type ModuleEntry struct {
	Name, Path string
}

type TemplateContext struct {
	ForFrontend bool
	Mods []ModuleEntry
}

const evalJs = `{{ range .Mods }}import * as {{ .Name }} from './{{ .Path }}';
{{ end }}
async function evalWith(cmd, me, value) {
    // the eval'd code can see the variables:
	//   * 'me' - the component which invoked the cmd
	//   * 'value' - the result of a previous command, usually run in the backend
	//   * all frontend JS modules by name.
    return await eval(cmd);
}

{{ if .ForFrontend }}window{{ else }}global{{ end }}.evalWith = evalWith;
`

func WriteEvalJs(mods []string, outfile string, for_frontend bool) error {
	writer, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer writer.Close()
	return WriteEvalJsWith(writer, mods, for_frontend)
}

func WriteEvalJsWith(writer io.Writer, mods []string, for_frontend bool) error {
	tpl := template.Must(template.New("eval.js").Parse(evalJs))
	return tpl.Execute(writer, createContext(mods, for_frontend))
}

func createContext(mods []string, for_frontend bool) TemplateContext {
	var entries = make([]ModuleEntry, len(mods))
	for i, mod := range mods {
		n := ChangExtension(path.Base(mod), "")
		p := ChangExtension(mod, ".js")
		entries[i] = ModuleEntry{Name: n, Path: p}
	}
	return TemplateContext{ForFrontend: for_frontend, Mods: entries}
}
