package src

import (
	"log"
	"os"
)

/// CopyFile copies input to output, panics on error.
func CopyFile(input, output string) {
	in, err := os.Open(input)
	if err != nil {
		log.Fatal("openFile", err)
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		log.Fatal("createFile", err)
	}
	defer out.Close()
	_, err = out.ReadFrom(in)
	if err != nil {
		log.Fatal("copyFile", err)
	}
}
