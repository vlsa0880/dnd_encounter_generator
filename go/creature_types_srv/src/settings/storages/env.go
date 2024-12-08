package storages

import (
	"context"
	"creature_types_srv/src/settings"

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
