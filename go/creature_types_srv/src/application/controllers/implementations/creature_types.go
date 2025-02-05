package controllers

import (
	"context"
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type CreatureTypes struct {
	creatureTypesDB dbinterfaces.CreatureTypesRepository
}

func NewGetCreatureTypes(creatureTypesDb dbinterfaces.CreatureTypesRepository) (*CreatureTypes, error) {
	if creatureTypesDb == nil {
		return nil, fmt.Errorf("bad db")
	}

	return &CreatureTypes{
		creatureTypesDB: creatureTypesDb,
	}, nil
}

func (controller *CreatureTypes) Get(ctx context.Context) (entities.CreatureTypes, error) {
	return controller.creatureTypesDB.GetCreatureTypes(ctx)
}

func (controller *CreatureTypes) Set(ctx context.Context, creatureTypes *entities.CreatureTypes) error {
	return controller.creatureTypesDB.UploadCreatureTypes(ctx, creatureTypes)
}
