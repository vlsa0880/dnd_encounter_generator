package mocks

import (
	"context"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/mocks"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type CreatureTypesValid struct {
	db dbinterfaces.CreatureTypesRepository
}

func NewCreatureTypesValid() *CreatureTypesValid {
	return &CreatureTypesValid{db: mocks.New()}
}

func (controller *CreatureTypesValid) Get(ctx context.Context) (entities.CreatureTypes, error) {
	return controller.db.GetCreatureTypes(ctx)
}

func (controller *CreatureTypesValid) Set(ctx context.Context, creatureTypes *entities.CreatureTypes) error {
	return controller.db.UploadCreatureTypes(ctx, creatureTypes)
}
