package data_handlers_interfaces

import (
	"context"
	creature_types "creature_types_srv/src/creature_types/data"
)

type DataHandler interface {
	GetCreatureTypes(ctx context.Context) creature_types.Types
}
