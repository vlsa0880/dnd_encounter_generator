package implementations

import (
	"fmt"
	"sync"

	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/storages"
)

var (
	instance settings.ISettingsManager
	once     sync.Once
)

func GetInstance() settings.ISettingsManager {
	once.Do(func() {
		manager := storages.CreateSettingsManager()
		var ok bool
		if instance, ok = manager.(settings.ISettingsManager); !ok {
			panic("Wrong interface created by factory")
		}
		err := instance.Init()
		if err != nil {
			panic(fmt.Sprintf("Cannot init settings manager: %v", err))
		}
		if err = instance.Load(); err != nil {
			panic(fmt.Sprintf("Cannot load settings from settings manager: %v", err))
		}
	})
	return instance
}
