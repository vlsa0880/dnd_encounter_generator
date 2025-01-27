package grpc

import (
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCRouter struct {
	Server *grpc.Server
	config Config
}

type Config struct {
	GRPC struct {
		Address string
		Network string
	}
}

func (router *GRPCRouter) Run() error {
	if router.Server == nil {
		return fmt.Errorf("bad grpc server")
	}
	listener, err := net.Listen(
		router.config.GRPC.Network,
		router.config.GRPC.Address,
	)
	if err != nil {
		return fmt.Errorf("can't create net listener: %s", err)
	}

	logger.GetInstance().Info(
		"grpc server started",
		zap.String("addr", listener.Addr().String()),
	)

	if err := router.Server.Serve(listener); err != nil {
		return fmt.Errorf("can't create start grpc server: %w", err)
	}
	return nil
}

func NewServer(settingsLoader settings.SettingsLoader) (*GRPCRouter, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("Bad settings loader")
	}

	var router GRPCRouter
	if err := settingsLoader.Load(&router.config); err != nil {
		return nil, fmt.Errorf("Can't load grpc server settings: %w", err)
	}

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
			logging.StartCall,
			logging.FinishCall,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.GetInstance().Error(
				"Recovered from panic",
			)

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	log_interceptor := logging.UnaryServerInterceptor(
		interceptorLogger(logger.GetInstance()),
		loggingOpts...)

	router.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			log_interceptor,
		),
	)

	return &router, nil
}

func (router *GRPCRouter) Stop() error {
	if router.Server == nil {
		return fmt.Errorf("Can't stop unconstructed")
	}
	logger.GetInstance().Info(
		"stopping grpc server",
		zap.String("address", router.config.GRPC.Address),
	)
	router.Server.GracefulStop()
	return nil
}

func interceptorLogger(zap_log *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		zaped_fields := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				zaped_fields = append(zaped_fields, zap.String(key.(string), v))
			case int:
				zaped_fields = append(zaped_fields, zap.Int(key.(string), v))
			case bool:
				zaped_fields = append(zaped_fields, zap.Bool(key.(string), v))
			default:
				zaped_fields = append(zaped_fields, zap.Any(key.(string), v))
			}
		}

		logger := zap_log.WithOptions(zap.AddCallerSkip(1)).With(zaped_fields...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
