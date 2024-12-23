package grpc

import (
	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/services"
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
	server *grpc.Server
	config Config
}

type Config struct {
	GRPC struct {
		Address string
		Network string
	}
}

func (router *GRPCRouter) Run() {
	if router.server == nil {
		panic("bad grpc server")
	}
	listener, err := net.Listen(
		router.config.GRPC.Network,
		router.config.GRPC.Address,
	)
	if err != nil {
		err_msg := fmt.Errorf("can't create net listener: %s", err)
		panic(err_msg)
	}

	logger.GetInstance().Info(
		"grpc server started",
		zap.String("addr", listener.Addr().String()),
	)

	if err := router.server.Serve(listener); err != nil {
		err_msg := fmt.Errorf("can't create start grpc server: %s", err)
		panic(err_msg)
	}
}

func New(settingsLoader settings.SettingsLoader) *GRPCRouter {
	if settingsLoader == nil {
		panic("Bad settings loader")
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

	router := GRPCRouter{}
	router.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			log_interceptor,
		),
	)

	if err := settingsLoader.Load(&router.config); err != nil {
		err_msg := fmt.Errorf("Can't load grpc server settings: %s", err)
		panic(err_msg)
	}

	return &router
}

func (router *GRPCRouter) SetupDataHandler(data_handler data_handlers_interfaces.DataHandler) error {
	if data_handler == nil {
		panic("Bad data handler")
	}
	services.New(router.server, data_handler)
	return nil
}

func (router *GRPCRouter) Stop() {
	if router.server == nil {
		return
	}
	logger.GetInstance().Info(
		"stopping grpc server",
		zap.String("address", router.config.GRPC.Address),
	)
	router.server.GracefulStop()
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
