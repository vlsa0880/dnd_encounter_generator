package logger

import (
	settings_data "creature_types_srv/src/settings"
	settings_manager "creature_types_srv/src/settings/implementations"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// TODO: get rid of strong cohesion
// TODO: add interface for logger settings for loose coupling
func GetInstance() *zap.Logger {
	once.Do(func() {
		instance = setupLogger()
		instance.Info(
			"logger inited",
			zap.String("settings", settings_manager.GetInstance().GetSettings().PrettyJson()),
		)
	})
	return instance
}

func setupLogger() *zap.Logger {
	var log *zap.Logger

	var err error
	switch settings_manager.GetInstance().GetSettings().Global.EnvType {
	case settings_data.EnvLocal:
		log, err = zap.NewDevelopment()
	case settings_data.EnvDev:
		log, err = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "json",
			OutputPaths: []string{"stdout", "/tmp/logs"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "message",
				LevelKey:     "level",
				TimeKey:      "time",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				CallerKey:    "caller",
				EncodeCaller: zapcore.FullCallerEncoder,
			},
		}.Build()
	case settings_data.EnvProd:
		log, err = zap.NewProduction()
	default:
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Can't create logger instance")
	}

	return log
}
