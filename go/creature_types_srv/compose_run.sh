docker image prune
docker builder prune
export UID=$(id -u)
export GID=$(id -g)
docker compose  \
    -f docker/compose/webapp.yaml \
    -f docker/compose/elastic_server.yaml \
    -f docker/compose/kafka_server.yaml \
    -f docker/compose/postgres.yaml \
    -f docker/compose/prometheus.yaml \
    up --build --force-recreate --remove-orphans --abort-on-container-exit
