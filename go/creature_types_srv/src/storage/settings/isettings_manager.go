package settings

import (
	"fmt"
	"sync"
)

var (
	instance ISettingsManager
	once     sync.Once
)

func GetInstance() ISettingsManager {
	once.Do(func() {
		instance = CreateSettingsManager()
		if instance == nil {
			panic("Cannot create settings manager")
		}
		err := instance.Init()
		if err != nil {
			panic(fmt.Sprintf("Cannot init settings manager %v", err))
		}
	})
	fmt.Printf("Settings inited: %v\n", instance.GetSettings().GetSlogGroup().String())
	return instance
}

type ISettingsManager interface {
	Init() error
	GetSettings() *SettingsData
}

type BaseSettingsManager struct {
	settings *SettingsData
}

func (base_settings *BaseSettingsManager) GetSettings() *SettingsData {
	return base_settings.settings
}
