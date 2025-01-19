mkdir -p bin && go build -o bin/client src/main.go && godotenv -f ../../../creature_types_srv/local/env/.env ./bin/client
