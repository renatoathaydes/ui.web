package src

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
)

type HasLogger interface {
	Log() *slog.Logger
}

// / CopyFile copies input to output, panics on error.
func CopyFile(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.ReadFrom(in)
	return err
}

func copyIndexHtml(logger *slog.Logger, input, output, dir string) error {
	var data struct {
		ImportMap string
	}
	b := make([]byte, 0, 1024)
	w := bytes.NewBuffer(b)
	err := GenerateImportMaps(logger, path.Join(dir, "node_modules"), w)
	if err != nil {
		return fmt.Errorf("error generating importmaps: %s", err)
	}
	logger.Debug("Generated importmap from node_modules", "dir", dir)
	data.ImportMap = w.String()
	tmpl, err := template.New("index.html").ParseFiles(input)
	if err != nil {
		return fmt.Errorf("could not parse index.html template: %s", err)
	}
	writer, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open template output file: %s", err)
	}
	defer writer.Close()
	err = tmpl.Execute(writer, &data)
	if err != nil {
		return fmt.Errorf("could not write parsed index.html template: %s", err)
	}
	return nil
}

func CopyAssetsDir(logger *slog.Logger, dir, destination string) error {
	source := path.Join(dir, "assets")
	var err error = filepath.Walk(source, func(p string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(p, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), os.ModePerm)
		} else if info.Mode().IsRegular() {
			var err error
			if relPath == "/index.html" {
				err = copyIndexHtml(logger, p, filepath.Join(destination, relPath), dir)
			} else {
				logger.Debug("Copying file", "from", relPath, "to", destination)
				err = CopyFile(p, filepath.Join(destination, relPath))
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func ChangExtension(file, newExt string) string {
	ext := path.Ext(file)
	if ext == newExt {
		return file
	}
	base := file[0 : len(file)-len(ext)]
	return base + newExt
}

// Start watching a directory.
//
// The trigger function is called on any detected changes.
// If the function returns true, the watcher continues, otherwise
// it will abort.
// This function returns immediately as the watcher runs async.
func WatchAsync[T HasLogger](dir string, ctx T, trigger func(T) bool, onExit func()) error {
	logger := ctx.Log()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	watcher.Add(dir)
	go func(w *fsnotify.Watcher, ctx T, trigger func(T) bool, onExit func()) {
		defer w.Close()
	loop:
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					break loop
				}
				// ignore as per docs, not useful:
				// https://github.com/fsnotify/fsnotify?tab=readme-ov-file#faq
				if !event.Has(fsnotify.Chmod) && isAtDir(event.Name, dir) {
					logger.Debug("File system change detected", "path", event.Name, "op", event.Op.String())
					keepGoing := trigger(ctx)
					if !keepGoing {
						break loop
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Warn("File watcher error", "error", err)
			}
		}
		onExit()
	}(watcher, ctx, trigger, onExit)

	return nil
}

// isAtDir tells whether a file is exactly at dir.
// It returns false if file is in a sub-directory of dir
// or if file itself is a directory.
func isAtDir(file, dir string) bool {
	d := path.Dir(file)
	if d != dir {
		return false
	}
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}
