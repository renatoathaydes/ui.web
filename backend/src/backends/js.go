package backends

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type CommandResponse struct {
	Value interface{} `json:value`
	Error *string     `json:error`
}

func RunJsCommand(cmd string) (*CommandResponse, error) {
	cmdBody := strings.NewReader(cmd)
	res, err := http.Post("http://localhost:8001/command", "text/plain", cmdBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	textBody := string(body)
	response := CommandResponse{}
	decoder := json.NewDecoder(strings.NewReader(textBody))
	jsErr := decoder.Decode(&response)
	if jsErr != nil {
		return nil, jsErr
	}
	return &response, nil
}
