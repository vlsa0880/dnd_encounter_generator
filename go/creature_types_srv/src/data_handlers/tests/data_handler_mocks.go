package data_handlers_tests

import (
	"context"
	"fmt"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	creature_types_tests "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/tests"
)

type DataHandlerMock struct {
	MockManager creature_types_tests.CreatureTypeManagerMockValid
}

func New() *DataHandlerMock {
	mock := DataHandlerMock{}
	if mock.MockManager.Init() != nil {
		return nil
	}
	return &mock
}

func (mock *DataHandlerMock) GetCreatureTypes(ctx context.Context) (*creature_types.Data, error) {
	return mock.MockManager.GetData()
}

func (handler *DataHandlerMock) PutCreatureTypes(ctx context.Context, types *creature_types.Data) error {
	return fmt.Errorf("not implemented")
}
