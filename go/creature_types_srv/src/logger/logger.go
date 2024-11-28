package logger

import (
	"dnd_encounter_builder/src/storage/settings"
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

func GetInstance() *slog.Logger {
	once.Do(func() {
		instance = setupLogger()
	})
	instance.Info(
		"logger created",
		settings.GetInstance().GetSettings().GetSlogGroup(),
	)
	return instance
}

func setupLogger() *slog.Logger {
	var log *slog.Logger

	switch settings.GetInstance().GetSettings().Logger.EnvType {
	case settings.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case settings.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case settings.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo},
			))
	default:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			))
	}

	return log
}
