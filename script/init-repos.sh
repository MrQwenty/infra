#!/bin/bash


PRINCIPAL_DIR="influenza-workspace"


REPOS=(
    "https://orbyta-isi-foundation@dev.azure.com/orbyta-isi-foundation/influenzanet/_git/user-management-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/study-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/messaging-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/logging-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/infra"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/api-gateway"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/participant-webapp"
)


mkdir -p "$PRINCIPAL_DIR"
cd "$PRINCIPAL_DIR" || exit 1

for REPO in "${REPOS[@]}"; do
    REPO_NAME=$(basename -s .git "$REPO")
    if [ -d "$REPO_NAME" ]; then
        echo "Repository $REPO_NAME already exists, skipping..."
    else
        git clone "$REPO" "$REPO_NAME" && cd "$REPO_NAME"
        if [[ "$REPO_NAME" == *infra* ]]; then
            git checkout develop
        else
            git checkout master
        fi
        cd ..
    fi
done


if ! docker info &>/dev/null; then
    echo "Docker is not running. Please start Docker and rerun the script."
    exit 1
fi

echo "Docker is running. Proceeding with Docker Compose in infra repositories."


cd infra

cd fe_config


cp ".env.local" "../../participant-webapp/"
echo "Copied .env.local to participant-webapp."

cp -r "public" "../../participant-webapp/"
echo "Copied public dir to participant-webapp."

cp "setupProxy.js" "../../participant-webapp/src/"
echo "Copied setupProxy.js to participant-webapp."

rm "../../participant-webapp/src/configs"
ln -s  "../../participant-webapp/public/assets/configs" "../../participant-webapp/src/configs"
echo "Symlinked to participant-webapp/public."

cd ..

if [ -f "docker-compose.yml" ]; then
            echo "Running Docker Compose in infra..."
            docker-compose up -d
else
            echo "No docker-compose.yml found in infra, skipping..."
fi

echo "InfluenzaNet is ready to GO 🚀🚀🚀🚀🚀🚀"

cd ../participant-webapp
echo "/src/configs" >> .gitignore
echo "/src/configs/*" >> .gitignore
