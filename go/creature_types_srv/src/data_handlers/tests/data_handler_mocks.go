package data_handlers_tests

import (
	"context"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	creature_types_tests "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/tests"
)

type DataHandlerMock struct {
	MockManager creature_types_tests.CreatureTypeManagerMockValid
}

func (mock *DataHandlerMock) GetCreatureTypes(ctx context.Context) creature_types.Types {
	return mock.MockManager.GetData()
}
