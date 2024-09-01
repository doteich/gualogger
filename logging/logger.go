package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(lvl string) {

	var sLvl slog.Level

	switch lvl {
	case "DEBUG":
		sLvl = slog.LevelDebug
	case "WARN":
		sLvl = slog.LevelWarn
	case "ERROR":
		sLvl = slog.LevelError
	default:
		sLvl = slog.LevelInfo
	}

	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: sLvl}))

}
