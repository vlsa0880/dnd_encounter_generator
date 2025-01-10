package settings

type SettingsLoader interface {
	Load(target interface{}) error
}
