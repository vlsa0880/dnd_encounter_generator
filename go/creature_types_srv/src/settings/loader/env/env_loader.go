package settings

import (
	"github.com/vrischmann/envconfig"
)

type EnvLoader struct {
}

func New() *EnvLoader {
	return &EnvLoader{}
}

func (loader *EnvLoader) Load(target interface{}) error {
	if err := envconfig.Init(target); err != nil {
		return err
	}
	return nil
}
