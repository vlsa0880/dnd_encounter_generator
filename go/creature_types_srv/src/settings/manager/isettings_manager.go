package settings

type ISettingsManager interface {
	Init() error
	Load() error
	GetSettings() *Data
}
