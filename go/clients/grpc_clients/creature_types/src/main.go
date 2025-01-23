package main

import (
	"context"
	"fmt"
	"time"

	irouter "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/grpc"
	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/grpc/gen"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initSettingsLoader() *settings.EnvLoader {
	loader := settings.New()
	if loader == nil {
		panic("bad settings loader")
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)
	return loader
}

func getServerConfig(loader *settings.EnvLoader) *irouter.Config {
	var config irouter.Config
	if err := loader.Load(&config); err != nil {
		panic(fmt.Errorf("Can't load grpc server settings: %s", err))
	}
	return &config
}

type Config struct {
	GRPC struct {
		ClientTimeout time.Duration
	}
}

func getClientConfig(loader *settings.EnvLoader) *Config {
	var config Config
	if err := loader.Load(&config); err != nil {
		panic(fmt.Errorf("Can't load grpc client settings: %s", err))
	}
	return &config
}

func createConnection(loader *settings.EnvLoader) *grpc.ClientConn {
	config := getServerConfig(loader)
	grpcConnection, err := grpc.NewClient(
		config.GRPC.Address,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	if err != nil {
		panic("Can't create grpc connector")
	}
	return grpcConnection
}

func main() {
	settingsLoader := initSettingsLoader()
	clientConfig := getClientConfig(settingsLoader)
	logger.GetInstance().Info(
		"client config created",
		zap.String("config", fmt.Sprintf("%v", clientConfig)),
	)

	conn := createConnection(settingsLoader)
	defer conn.Close()

	client := gen.NewDataHandlerClient(conn)

	req := gen.GetRequest{}
	ctx, cancel := context.WithTimeout(
		context.Background(),
		clientConfig.GRPC.ClientTimeout,
	)
	defer cancel()
	var resp *gen.GetResponse
	watcherChannel := make(chan struct{})
	go func() {
		var err error
		if resp, err = client.GetCreatureTypes(ctx, &req); err != nil {
			err_msg := fmt.Sprintf("Can't get creature types by grpc: %s", err)
			panic(err_msg)
		}
		watcherChannel <- struct{}{}
	}()

	select {
	case <-watcherChannel:
		break
	case <-ctx.Done():
		panic("get creature type stopped by timeout")
	}

	fmt.Println(resp)
}
