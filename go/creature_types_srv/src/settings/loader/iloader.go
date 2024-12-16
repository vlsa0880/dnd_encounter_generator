package settings

type ISettingsLoader interface {
	Load(target interface{}) error
}
