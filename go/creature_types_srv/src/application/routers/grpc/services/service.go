package services

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/grpc/gen"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type serverAPI struct {
	gen.UnimplementedDataHandlerServer
	creatureTypesDB dbinterfaces.CreatureTypesRepository
}

func Run(gRPCServer *grpc.Server, creatureTypesDB dbinterfaces.CreatureTypesRepository) {
	gen.RegisterDataHandlerServer(gRPCServer, &serverAPI{creatureTypesDB: creatureTypesDB})
}

func (server *serverAPI) GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	creatureTypes, err := server.creatureTypesDB.GetCreatureTypes(ctx)
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
