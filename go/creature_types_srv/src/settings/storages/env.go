package storages

import (
	"context"

	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager"

	"github.com/sethvargo/go-envconfig"
)

type Env struct {
	BaseSettingsManager
}

func (env *Env) Init() error {
	return nil
}

func (env *Env) Load() error {
	env.settings = settings.NewData()
	if err := envconfig.Process(context.Background(), env.settings); err != nil {
		return err
	}
	return nil
}
