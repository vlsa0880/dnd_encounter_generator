package logger

import (
	settings_data "creature_types_srv/src/settings"
	settings_manager "creature_types_srv/src/settings/implementations"
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
		settings_manager.GetInstance().GetSettings().GetSlogGroup(),
	)
	return instance
}

func setupLogger() *slog.Logger {
	var log *slog.Logger

	switch settings_manager.GetInstance().GetSettings().Global.EnvType {
	case settings_data.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case settings_data.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case settings_data.EnvProd:
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
