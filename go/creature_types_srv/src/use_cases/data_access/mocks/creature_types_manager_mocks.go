package mocks

import (
	"context"
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
)

type CreatureTypeManagerMockValid struct {
	types entities.CreatureTypes
}

func New() *CreatureTypeManagerMockValid {
	return &CreatureTypeManagerMockValid{
		types: entities.CreatureTypes{
			Types: []entities.CreatureType{
				{Name: "тест1"},
				{Name: "тест2"},
			},
		},
	}
}

func (mock *CreatureTypeManagerMockValid) GetCreatureTypes(ctx context.Context) (entities.CreatureTypes, error) {
	return mock.types, nil
}

func (mock *CreatureTypeManagerMockValid) UploadCreatureTypes(ctx context.Context, types *entities.CreatureTypes) error {
	if types == nil {
		return fmt.Errorf("bad creature types")
	} else if len(types.Types) == 0 {
		return fmt.Errorf("empty types")
	}
	mock.types = *types
	return nil
}
