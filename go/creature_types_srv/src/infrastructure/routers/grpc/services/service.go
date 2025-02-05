package services

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/grpc/gen"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
)

type serverAPI struct {
	gen.UnimplementedDataHandlerServer
	controller icontrollers.CreatureTypes
}

func Run(gRPCServer *grpc.Server, controller icontrollers.CreatureTypes) {
	gen.RegisterDataHandlerServer(gRPCServer, &serverAPI{controller: controller})
}

func (server *serverAPI) GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	creatureTypes, err := server.controller.Get(ctx)
	if err != nil {
		logger.GetInstance().Error(
			"can't get creature types",
			zap.Error(err),
		)
	}
	resp := gen.GetResponse{
		CreatureTypes: make([]*gen.CreatureType, 0, len(creatureTypes.Types)),
	}
	for _, creature_type := range creatureTypes.Types {
		resp.CreatureTypes = append(
			resp.CreatureTypes,
			&gen.CreatureType{Name: creature_type.Name})
	}
	return &resp, nil
}
