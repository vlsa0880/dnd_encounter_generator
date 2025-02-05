package icontrollers

import (
	"context"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
)

type CreatureTypes interface {
	Get(ctx context.Context) (entities.CreatureTypes, error)
	Set(ctx context.Context, creatureTypes *entities.CreatureTypes) error
}
