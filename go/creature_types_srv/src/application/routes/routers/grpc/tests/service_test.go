package tests

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gen "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/grpc/gen"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/grpc/services"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/mocks"
)

func TestGetCreatureTypes_GRPCCall(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer listener.Close()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	dataHandler := mocks.New()
	assert.NotNil(t, dataHandler)
	expectedData, err := dataHandler.GetCreatureTypes(context.Background())
	assert.NoError(t, err)
	services.New(grpcServer, dataHandler)
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			t.Error(
				fmt.Printf("grpc server return an error: %s", err),
			)
		}
	}()

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
	validFinishChannel := make(chan struct{})
	go func() {
		resp, err = client.GetCreatureTypes(ctx, &req)
		validFinishChannel <- struct{}{}
	}()

	select {
	case <-validFinishChannel:
		break
	case <-ctx.Done():
		t.Error("get creature type stopped by timeout")
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	for index, creature_type_resp := range resp.CreatureTypes {
		assert.Equal(t, creature_type_resp.Name, expectedData.Types[index].Name)
	}
}
