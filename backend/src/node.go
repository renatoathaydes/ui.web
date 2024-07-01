package src

import (
	"log/slog"
	"os/exec"
	"path"
)

func StartNode(modsDir string, logger *slog.Logger) {
	cmd := exec.Command("node", "--watch", path.Join(modsDir, "out", "startup.js"))
	node_logger := logger.WithGroup("node")
	stdout := NewWriter(func(line string) {
		node_logger.Info(line)
	})
	stderr := NewWriter(func(line string) {
		node_logger.Warn(line)
	})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		logger.Error("backend-js node process could not be started", "error", err)
	}

	go func() {
		defer stdout.Close()
		defer stderr.Close()
		err := cmd.Wait()
		if err != nil {
			logger.Error("backend-js node process died", "error", err)
		}
	}()
}
