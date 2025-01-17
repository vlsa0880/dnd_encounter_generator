mkdir -p bin && go build -o bin/client src/main.go && godotenv -f ../../../creature_types_srv/docker/env/webapp.env ./bin/client
