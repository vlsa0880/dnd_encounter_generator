# Creatures service

## Running with docker compose
### Rebuild only on src/ files change
```sh
./compose_run_watch.sh
```
### Forced full rebuild
```sh
export UID=$(id -u)
export GID=$(id -g)
docker compose  \
    -f docker/compose/webapp.yaml \
    -f docker/compose/elastic_server.yaml \
    -f docker/compose/kafka_server.yaml \
    -f docker/compose/postgres.yaml \
    -f docker/compose/prometheus.yaml \
    up --build --force-recreate --remove-orphans
```
or
```sh
./compose_run.sh
```

## Running localy
### Postgres db launch
```sh
./run_psql_docker.sh
```

### Service build
```sh
source local/env/.env
export $(cut -d= -f1 local/env/.env)
cd go/creature_types_srv/
mkdir -p bin && go build -C bin/ ../src/main.go && ./bin/main
```

## Request available creature types
### Docker
```sh
curl "0.0.0.0:8088/creatures_data/types"
```

### Localy
```sh
curl "localhost:8088/creatures_data/types"
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
