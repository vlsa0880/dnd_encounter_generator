# Creatures service

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

## Running by docker compose
### Rebuild only on src/ files change
```sh
./compose_run_watch.sh
```
### Forced full rebuild
```sh
docker compose -f compose.yaml up --build --force-recreate --remove-orphans
```
or
```sh
./compose_run.sh
```

## Request available creature types
### When running with docker
```sh
curl "0.0.0.0:8088/creature_data/types"
```

### When running localy
```sh
curl "localhost:8088/creature_data/types"
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
cd ../grpc_clients/creature_types/
```

After running docker compose & grpc server started:
```sh
mkdir -p bin && go build -o bin/client src/main.go && ./run_client_docker.sh
```

Or localy, after running psql & get_creature_srv:
```sh
mkdir -p bin && go build -o bin/client src/main.go && ./run_client_localy.sh
```
