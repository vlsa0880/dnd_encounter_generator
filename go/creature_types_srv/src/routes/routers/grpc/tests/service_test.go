package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	data_handlers_tests "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/tests"
	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/gen"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/services"
)

func TestGetCreatureTypes_GRPCCall(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer listener.Close()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	data_handler := data_handlers_tests.DataHandlerMock{}
	expected_data := data_handler.GetCreatureTypes(context.Background())

	services.New(grpcServer, &data_handler)
	go grpcServer.Serve(listener)

	grpcConnect, err := grpc.NewClient(
		listener.Addr().String(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	assert.NoError(t, err)
	defer grpcConnect.Close()

	client := gen.NewDataHandlerClient(grpcConnect)

	req := gen.GetRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var resp *gen.GetResponse
	watcher_channel := make(chan struct{})
	go func() {
		resp, err = client.GetCreatureTypes(ctx, &req)
		watcher_channel <- struct{}{}
	}()

	select {
	case <-watcher_channel:
		break
	case <-ctx.Done():
		t.Error("get creature type stopped by timeout")
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	for index, creature_type_resp := range resp.CreatureTypes {
		assert.Equal(t, creature_type_resp.RU, expected_data[index].RU)
	}
}
