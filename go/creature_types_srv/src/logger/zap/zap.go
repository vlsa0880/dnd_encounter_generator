package zap_logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/kafka"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	envdata "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/data"

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
		panic("Logger wasn't initialised")
	}
	return instance
}

func InitLogger(settingsLoader settings.SettingsLoader) {
	once.Do(func() {
		config := Config{}
		if err := settingsLoader.Load(&config); err != nil {
			err_msg := fmt.Sprintf("Can't load logger config: %s", err)
			panic(err_msg)
		}
		instance = setupLogger(&config, settingsLoader)
		instance.Info("logger inited")
	})
}

func setupLogger(config *Config, settingsLoader settings.SettingsLoader) *zap.Logger {
	var log *zap.Logger

	var err error
	switch config.Global.EnvType {
	case envdata.Local.String():
		log, err = zap.NewDevelopment()
	case envdata.Dev.String():
		log, err = newDev(settingsLoader)
	case envdata.Prod.String():
		log, err = zap.NewProduction()
	default:
		panic("No environment type specified")
	}

	if err != nil {
		panic(fmt.Errorf("Can't create logger instance: %w", err))
	}

	return log
}

func newDev(settingsLoader settings.SettingsLoader) (*zap.Logger, error) {
	kafkaSyncer, err := kafka.New(settingsLoader)
	if err != nil {
		return nil, fmt.Errorf("Can't create kafka log syncer: %w", err)
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(kafkaSyncer),
			zapcore.DebugLevel,
		),
	)
	return zap.New(core, zap.AddCaller()), nil
}
