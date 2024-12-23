module creature_types_grpc_client

go 1.23.2

require google.golang.org/grpc v1.69.2

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0 // indirect
	github.com/vrischmann/envconfig v1.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

require (
	github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv v0.0.0
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241219192143-6b3ec007d9bb // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)

replace github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv v0.0.0 => ../../creature_types_srv
