package interfaces

import (
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/interfaces"
)

type Consumer interface {
	interfaces.Runner
	MsgReceiver
}
