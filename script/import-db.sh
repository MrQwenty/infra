#!/bin/bash

DUMP_FOLDER="$(realpath "$(dirname "$0")/../dmp_influenza/mongo.dump")"
MONGO_CONTAINER_NAME="mongodb"
MONGO_USERNAME="admin"
MONGO_PASSWORD="root"
DOCKER_COMPOSE_FILE="$(realpath "$(dirname "$0")/../docker-compose.yml")"

function wait_for_mongodb() {
    echo "Waiting for mongodb to be ready..."
    until docker exec "$MONGO_CONTAINER_NAME" mongosh --eval "db.runCommand({ ping: 1 })" --username "$MONGO_USERNAME" --password "$MONGO_PASSWORD" --quiet &>/dev/null; do
        echo "MongoDB is not ready yet. Retrying in 5 seconds..."
        sleep 5
    done
    echo "Mongodb is ready"
}

if ! docker ps --format '{{.Names}}' | grep -q "$MONGO_CONTAINER_NAME"; then
    echo "Error: MongoDB container '$MONGO_CONTAINER_NAME' is not running"
    exit 1
fi

wait_for_mongodb

# Proceed with data import
echo "Copying MongoDB dump to container..."
docker cp "$DUMP_FOLDER" "$MONGO_CONTAINER_NAME":/data

echo "Restoring dump..."
docker exec -i "$MONGO_CONTAINER_NAME" mongorestore --archive=data/mongo.dump --username "$MONGO_USERNAME" --password "$MONGO_PASSWORD"

echo "Data import completed"

# Restart all containers to refresh services
echo "Restarting all services in Docker Compose..."
docker-compose -f "$DOCKER_COMPOSE_FILE" restart

echo "All services restarted successfully"
sleep 10
