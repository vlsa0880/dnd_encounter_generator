# Creatures service

## Running with docker compose
```sh
./compose_run.sh
```

## Request available creature types
### Docker
```sh
curl "0.0.0.0:8088/creatures_data/types"
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

## Tests run
```sh
go test -v ./src/...
```
