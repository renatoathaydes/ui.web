package backends

import (
	"io"
	"net/http"
	"strings"
)

func RunJsCommand(cmd string) (interface{}, error) {
	cmdBody := strings.NewReader(cmd)
	res, err := http.Post("http://localhost:8080/", "text/plain", cmdBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}
