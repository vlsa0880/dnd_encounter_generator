package interfaces

import (
	"context"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
)

type CreatureTypes interface {
	GetCreatureTypes(context.Context) (entities.CreatureTypes, error)
	UploadCreatureTypes(context.Context, *entities.CreatureTypes) error
}
