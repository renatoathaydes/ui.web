package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	p "path"
	"path/filepath"
)

func GenerateImportMaps(logger *slog.Logger, dir string, w io.Writer) error {
	_, err := w.Write([]byte("{\n  \"imports\": {"))
	if err != nil {
		return err
	}
	var nextDirs []string
	first := true
	nextDirs = append(nextDirs, dir)
outer:
	for {
		if len(nextDirs) == 0 {
			break outer
		}
		path := nextDirs[0]
		nextDirs = nextDirs[1:]

		entries, err := os.ReadDir(path)

		if err != nil {
			return err
		}
		relDir, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.Type().IsRegular() && entry.Name() == "package.json" {
				if first {
					first = false
					w.Write([]byte("\n"))
				} else {
					w.Write([]byte(",\n"))
				}
				err = writeImport(logger, relDir, path, w)
				if err != nil {
					return err
				}
				continue outer
			}
		}
		for _, entry := range entries {
			if entry.IsDir() {
				nextDirs = append(nextDirs, p.Join(path, entry.Name()))
			}
		}
	}
	_, err = w.Write([]byte("\n  }\n}\n"))
	return err
}

type nodePackage struct {
	Main    interface{} `json:"main"`
	Exports interface{} `json:"exports"`
}

// See for gory details: https://webpack.js.org/guides/package-exports/
func writeImport(logger *slog.Logger, relDir, dir string, w io.Writer) error {
	file, err := os.Open(p.Join(dir, "package.json"))
	if err != nil {
		return fmt.Errorf("cannot open %s package.json file: %s", relDir, err)
	}
	defer file.Close()
	_, err = w.Write([]byte("    "))
	if err != nil {
		return fmt.Errorf("cannot write to %v: %s", w, err)
	}
	j, err := json.Marshal(relDir)
	if err != nil {
		return fmt.Errorf("cannot marshall string: %s", err)
	}
	_, err = w.Write(j)
	if err != nil {
		return fmt.Errorf("cannot write to %v: %s", w, err)
	}
	_, err = w.Write([]byte(": "))
	if err != nil {
		return fmt.Errorf("cannot write to %v: %s", w, err)
	}
	pkg := nodePackage{}
	err = json.NewDecoder(file).Decode(&pkg)
	if err != nil {
		return fmt.Errorf("cannot decode JSON package.json at %s: %s", relDir, err)
	}
	if pkg.Exports == nil {
		err = writeMain(w, relDir, pkg.Main)
		if err != nil {
			return err
		}
	} else if exp, ok := pkg.Exports.(string); ok {
		err = writeEntryStr(w, relDir, exp)
		if err != nil {
			return err
		}
	} else if exp, ok := pkg.Exports.(map[string]interface{}); ok {
		err = writeImportForExport(logger, w, relDir, exp)
		if err != nil {
			return err
		}
	} else {
		logger.Warn("package exports not recognized", "pkg", relDir, "unsupported_type", pkg.Exports)
	}
	return nil
}

func writeMain(w io.Writer, relDir string, exp interface{}) error {
	if s, ok := exp.(string); ok {
		return writeEntryStr(w, relDir, s)
	}
	return fmt.Errorf("value of field 'main' in package '%s' is not a string: %v", relDir, exp)
}

func writeEntryStr(w io.Writer, relDir, exp string) error {
	b, err := json.Marshal("./" + p.Join("node_modules", relDir, exp))
	if err != nil {
		return fmt.Errorf("at %s, cannot marshall to %v: %s", relDir, w, err)
	}
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("at %s, cannot write to %v: %s", relDir, w, err)
	}
	return nil
}

func writeImportForExport(logger *slog.Logger, w io.Writer, relDir string, exp map[string]interface{}) error {
	// prefer the "import" field
	if imp, ok := exp["import"]; ok {
		if s, ok := imp.(string); ok {
			return writeEntryStr(w, relDir, s)
		} else {
			logger.Warn("Unrecognized 'import' selector", "pkg", relDir, "unsupported_type", imp)
		}
	} else {
		logger.Warn("Export missing 'import' selector", "pkg", relDir)
	}
	return nil
}
