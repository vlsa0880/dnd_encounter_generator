package settings

import (
	"log/slog"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

const (
	c_settings_storage_type_env     = "SETTINGS_STORAGE_TYPE"
	c_default_settings_storage_type = "postgres"
	c_storage_name_env              = "STORAGE_NAME"
	c_default_storage_name          = "postgres"
	c_storage_address_env           = "STORAGE_ADDRESS"
	c_default_storage_address       = "localhost"
)

type SettingsData struct {
	Http   `json:"http"`
	Kafka  `json:"kafka"`
	Logger `json:"logger"`
}

type Http struct {
	GetCreatureTypesEndpoint string `json:"get_creature_type_endpoint" default:"/creature_data/available_types"`
}

type Kafka struct {
	Servers                  string `json:"server" default:"localhost"`
	UpdateCreatureTypesTopic string `json:"update_creature_type_topic" default:"update_creature_types"`
}

type Logger struct {
	EnvType string `json:"env_type" default:"dev"`
}

func newSettingsData() *SettingsData {
	data := SettingsData{}
	data.Http.GetCreatureTypesEndpoint = "/creature_data/available_types"
	data.Kafka.Servers = "localhost"
	data.Kafka.UpdateCreatureTypesTopic = "update_creature_types"
	data.Logger.EnvType = "dev"
	return &data
}

func (settings *SettingsData) GetSlogGroup() slog.Attr {
	return slog.Group(
		"settings",
		slog.String("logger.env_type", settings.Logger.EnvType),
		slog.String("http.creature_type_ep", settings.Http.GetCreatureTypesEndpoint),
		slog.String("kafka.servers", settings.Kafka.Servers),
		slog.String("kafka.update_creature_type_topic", settings.Kafka.UpdateCreatureTypesTopic),
	)
}
