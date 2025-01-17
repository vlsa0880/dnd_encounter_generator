module github.com/vlsa0880/dnd_encounter_generator/go/clients/kafka_clients/creature_types

go 1.23.4

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/google/uuid v1.6.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv v0.0.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/vrischmann/envconfig v1.3.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv v0.0.0 => ../../../creature_types_srv
