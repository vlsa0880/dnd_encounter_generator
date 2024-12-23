package logger

import (
	"log/slog"
	"os"
	"sync"

	settings_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager"
	settings_manager "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager/implementations"
)

var (
	instance *slog.Logger
	once     sync.Once
)

func GetInstance() *slog.Logger {
	once.Do(func() {
		instance = setupLogger()
		instance.Info(
			"logger created",
			settings_manager.GetInstance().GetSettings().GetSlogGroup(),
		)
	})
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
