package zap_logger

import (
	"fmt"
	"sync"

	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	settings_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

type Config struct {
	Global struct {
		EnvType string
	}
}

func GetInstance() *zap.Logger {
	if instance == nil {
		panic("Logger doesn't initialised")
	}
	return instance
}

func InitLogger(settings_loader settings.SettingsLoader) {
	once.Do(func() {
		config := Config{}
		if err := settings_loader.Load(&config); err != nil {
			err_msg := fmt.Sprintf("Can't load logger config: %s", err)
			panic(err_msg)
		}
		instance = setupLogger(&config)
		instance.Info("logger inited")
	})
}

func setupLogger(config *Config) *zap.Logger {
	var log *zap.Logger

	var err error
	switch config.Global.EnvType {
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
		panic("No environment type specified")
	}

	if err != nil {
		panic("Can't create logger instance")
	}

	return log
}
