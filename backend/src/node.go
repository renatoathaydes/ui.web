package src

import (
	"log/slog"
	"os/exec"
	"path"
	"time"
)

// StartNode starts the node process with the --watch flag.
//
// That ensures that when the backend-js project is rebuilt, changes will be picked up
// and the node process automatically restarted.
//
// It returns a channel that sends `true` when the node process crashes, signalling that
// the process should be restarted, or `false` if node finishes without errors, as it's
// assumed that in this case, the user explicitly stopped the node process.
func StartNode(modsDir string, logger *slog.Logger) chan bool {
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

	result := make(chan bool)

	go func() {
		defer stdout.Close()
		defer stderr.Close()
		err := cmd.Wait()
		if err != nil {
			logger.Error("node process died, will restart in a few seconds", "error", err)
			time.Sleep(2 * time.Second)
			result <- true
		} else {
			logger.Info("node process exited successfully, will not restart it.")
			result <- false
		}
	}()

	return result
}
