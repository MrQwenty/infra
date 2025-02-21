#!/bin/bash

echo "Starting selective rebuild process for InfluenzaNet services..."

# Move up to the infra directory
cd ..

# Ensure Docker is running
if ! docker info &>/dev/null; then
    echo "Docker is not running. Please start Docker and rerun the script."
    exit 1
fi

# Ensure docker-compose.yml exists
if [ ! -f "docker-compose.yml" ]; then
    echo "Error: docker-compose.yml not found in infra/. Ensure you are running this script from infra/scripts/."
    exit 1
fi

# Get available services from docker-compose.yml
AVAILABLE_SERVICES=$(docker-compose config --services)

echo "Available services:"
echo "$AVAILABLE_SERVICES"

echo "Enter the services you want to rebuild (separated by space), or type 'all' to rebuild everything:"
read -r -a SELECTED_SERVICES

# Handle frontend configuration if participant-webapp is selected
if [[ " ${SELECTED_SERVICES[*]} " == *"participant-webapp"* || " ${SELECTED_SERVICES[*]} " == *"all"* ]]; then
    echo "Setting up frontend configuration..."

    cd fe_config || { echo "Error: fe_config directory not found."; exit 1; }

    cp ".env.local" "../../participant-webapp/"
    echo "Copied .env.local to participant-webapp."

    cp -r "public" "../../participant-webapp/"
    echo "Copied public dir to participant-webapp."

    cp "setupProxy.js" "../../participant-webapp/src/"
    echo "Copied setupProxy.js to participant-webapp."

    rm -rf "../../participant-webapp/src/configs"
    ln -s  "../../participant-webapp/public/assets/configs" "../../participant-webapp/src/configs"
    echo "Symlinked to participant-webapp/public."

    cd ..
fi

# Rebuild services
if [[ " ${SELECTED_SERVICES[*]} " == *"all"* ]]; then
    echo "Rebuilding all services..."
    docker-compose down
    docker-compose up --build -d
    # Execute import-db.sh if it exists
    if [ -f "./script/import-db.sh" ]; then
        echo "Running import-db.sh..."
        chmod +x ./script/import-db.sh
        ./script/import-db.sh
    else
        echo "Skipping database import, file not found."
        sleep 5
    fi

else
    echo "Stopping selected services..."
    docker-compose stop "${SELECTED_SERVICES[@]}"

    echo "Removing old containers..."
    docker-compose rm -f "${SELECTED_SERVICES[@]}"

    echo "Rebuilding selected services..."
    docker-compose up --build -d "${SELECTED_SERVICES[@]}"

    echo "Starting new container ..."
    docker-compose up -d "${SELECTED_SERVICES[@]}"
fi

echo "Selected services successfully rebuilt and restarted!"
sleep 5