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
```sh
docker compose -f compose.yaml up --build --force-recreate --remove-orphans
```
or
```sh
./compose_run.sh
```