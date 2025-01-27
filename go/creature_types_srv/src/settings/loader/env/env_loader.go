package settings

import (
	"github.com/vrischmann/envconfig"
)

type EnvLoader struct {
}

func New() (*EnvLoader, error) {
	return &EnvLoader{}, nil
}

func (loader *EnvLoader) Load(target interface{}) error {
	if err := envconfig.Init(target); err != nil {
		return err
	}
	return nil
}
