package src

import (
	"log/slog"
	"os/exec"
	"path"
	"time"
)

func StartNode(modsDir string, logger *slog.Logger) {
	cmd := exec.Command("node", "--watch", path.Join(modsDir, "out", "startup.js"))

	node_stdout := slog.New(logger.Handler().WithAttrs([]slog.Attr{
		{Key: "proc", Value: slog.StringValue("node")},
		{Key: "s", Value: slog.StringValue("stdout")},
	}))
	node_stderr := slog.New(logger.Handler().WithAttrs([]slog.Attr{
		{Key: "proc", Value: slog.StringValue("node")},
		{Key: "s", Value: slog.StringValue("stdout")},
	}))
	stdout := NewWriter(func(line string) {
		node_stdout.Info(line)
	})
	stderr := NewWriter(func(line string) {
		node_stderr.Info(line)
	})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		logger.Error("node process could not be started", "error", err)
	} else {
		logger.Debug("node process started")
	}

	go func() {
		defer stdout.Close()
		defer stderr.Close()
		err := cmd.Wait()
		if err != nil {
			logger.Error("node process died, will restart in a few seconds", "error", err)
			time.Sleep(2 * time.Second)
			StartNode(modsDir, logger)
		} else {
			logger.Info("node process exited successfully.")
		}
	}()
}
