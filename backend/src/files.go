package src

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

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

func CopyDir(source, destination string) error {
	var err error = filepath.Walk(source, func(p string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(p, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), os.ModePerm)
		} else if info.Mode().IsRegular() {
			err := CopyFile(filepath.Join(source, relPath), filepath.Join(destination, relPath))
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
func WatchAsync[T any](dir string, ctx T, trigger func(T) bool, onExit func()) error {
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
					log.Printf("File system change detected: %s\n", event)
					keepGoing := trigger(ctx)
					if !keepGoing {
						break loop
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
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
