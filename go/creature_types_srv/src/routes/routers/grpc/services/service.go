package services

import (
	"context"

	"google.golang.org/grpc"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/gen"
)

type serverAPI struct {
	gen.UnimplementedDataHandlerServer
	data_handler data_handlers_interfaces.DataHandler
}

type Getter interface {
	GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error)
}

func New(gRPCServer *grpc.Server, handler data_handlers_interfaces.DataHandler) {
	gen.RegisterDataHandlerServer(gRPCServer, &serverAPI{data_handler: handler})
}

func (server *serverAPI) GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	resp := gen.GetResponse{}
	for _, creature_type := range server.data_handler.GetCreatureTypes(ctx) {
		resp.CreatureTypes = append(
			resp.CreatureTypes,
			&gen.CreatureType{Name: creature_type.Name})
	}
	return &resp, nil
}
