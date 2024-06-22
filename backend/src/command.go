package src

import (
	"fmt"

	"ui.web/server/src/backends"
)

func runCommand(cmd, lang string) (interface{}, error) {
	if lang == "js" {
		v, e := backends.RunJsCommand(cmd)
		return v, e
	}
	return nil, fmt.Errorf("no backend found for language: '%s'", lang)
}
