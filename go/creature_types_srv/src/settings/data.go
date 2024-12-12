package settings

import (
	"encoding/json"
	"log/slog"

	"github.com/creasty/defaults"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Data struct {
	Http          `json:"http" env:", prefix=HTTP_"`
	CreatureTypes `json:"creature_types" env:", prefix=CREATURE_TYPES_"`
	Kafka         `json:"kafka" env:", prefix=KAFKA_"`
	Global        `json:"global" env:", prefix=GLOBAL_"`
}

type Http struct {
	Address string `json:"address" env:"ADDRESS, overwrite" default:"localhost"`
	Port    string `json:"port" env:"PORT, overwrite" default:"8088"`
}

type CreatureTypes struct {
	GetEndpoint string `json:"get_endpoint" env:"GET_ENDPOINT, overwrite" default:"/creature_data/types"`
	StorageType string `json:"storage_type" env:"STORAGE_TYPE, overwrite" default:"postgres"`
}

type Kafka struct {
	Servers                  string `json:"server" env:"SERVERS, overwrite" default:"localhost"`
	UpdateCreatureTypesTopic string `json:"update_creature_type_topic" env:"UPDATE_CREATURE_TYPES_TOPIC, overwrite" default:"update_creature_types"`
}

type Global struct {
	EnvType string `json:"env_type" env:"ENV_TYPE, overwrite" default:"dev"`
}

func NewData() *Data {
	data := &Data{}
	if err := defaults.Set(data); err != nil {
		panic(err)
	}
	return data
}

func (settings *Data) GetSlogGroup() slog.Attr {
	return slog.Group(
		"settings",
		slog.String("logger.env_type", settings.Global.EnvType),
		slog.String("http.creature_type_ep", settings.CreatureTypes.GetEndpoint),
		slog.String("creatures_types.storage_type", settings.CreatureTypes.StorageType),
		slog.String("kafka.servers", settings.Kafka.Servers),
		slog.String("kafka.update_creature_type_topic", settings.Kafka.UpdateCreatureTypesTopic),
	)
}

func (settings *Data) String() string {
	json, err := json.Marshal(settings)
	if err != nil {
		panic("Can't transform settings to string json")
	}
	return string(json)
}

func (settings *Data) PrettyJson() string {
	json, err := json.MarshalIndent(settings, "", "	")
	if err != nil {
		panic("Can't transform settings to pretty json")
	}
	return string(json)
}
