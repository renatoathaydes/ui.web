package logui

import (
	"log/slog"

	"github.com/fatih/color"
)

type logColors struct {
	debug *color.Color
	info  *color.Color
	warn  *color.Color
	error *color.Color
}

var colors logColors = logColors{
	debug: color.New(color.BgHiBlack, color.FgWhite),
	info:  color.New(color.BgBlack, color.FgHiWhite),
	warn:  color.New(color.BgYellow, color.FgBlack),
	error: color.New(color.BgRed, color.FgHiWhite),
}

func levelColor(level slog.Level) *color.Color {
	switch level {
	case slog.LevelDebug:
		return colors.debug
	case slog.LevelError:
		return colors.error
	case slog.LevelInfo:
		return colors.info
	case slog.LevelWarn:
		return colors.warn
	default:
		return colors.info
	}
}
