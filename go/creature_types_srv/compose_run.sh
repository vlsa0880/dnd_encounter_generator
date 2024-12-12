export UID=$(id -u)
export GID=$(id -g)
docker compose -f docker/compose.yaml up --build --force-recreate --remove-orphans
