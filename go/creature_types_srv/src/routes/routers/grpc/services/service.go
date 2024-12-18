package services

import (
	"context"

	"google.golang.org/grpc"

	"creature_types_srv/src/data_handlers/interfaces"
	gen "creature_types_srv/src/routes/routers/grpc/gen"
)

type serverAPI struct {
	gen.UnimplementedDataHandlerServer
	data_handler interfaces.DataHandler
}

type Getter interface {
	GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error)
}

func New(gRPCServer *grpc.Server, handler interfaces.DataHandler) {
	gen.RegisterDataHandlerServer(gRPCServer, &serverAPI{data_handler: handler})
}

func (server *serverAPI) GetCreatureTypes(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	resp := gen.GetResponse{}
	for _, creature_type := range server.data_handler.GetCreatureTypes() {
		resp.CreatureTypes = append(resp.CreatureTypes, creature_type.RU)
	}
	return &resp, nil
}
