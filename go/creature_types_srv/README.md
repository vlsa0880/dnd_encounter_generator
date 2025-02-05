# Creatures service

## Running with docker compose
```sh
./compose_run.sh
```

## Replace creature type with new ones
```sh
curl -X PUT -H "Content-Type: application/json" -d '{"types": ["тест1", "тест2"]}' "localhost:8088/creatures/types"
```

## Request available creature types
```sh
curl "localhost:8088/creatures/types"
```

## GRPC
### Recreate contracts
[Install grpc](https://grpc.io/docs/languages/go/quickstart/)

Recreate grpc go contracts
```sh
./gen_grpc.sh
```

### Run separated client
```sh
cd go/clients/grpc_clients/creature_types/
```

After running docker compose & grpc server started:
```sh
mkdir -p bin && go build -o bin/client src/main.go && ./run_client_docker.sh
```

### Same for kafka, except for folder
```sh
go/clients/kafka_clients/creature_types
```
WARN: for now test kafka client run correctly only after second run

## Tests run
```sh
go test -v ./src/...
```
