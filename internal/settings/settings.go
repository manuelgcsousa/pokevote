package settings

import (
	"log/slog"
	"os"
)

func Setup() {
	// Init logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}
