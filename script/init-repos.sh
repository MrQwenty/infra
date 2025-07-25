#!/bin/bash


PRINCIPAL_DIR="influenza-workspace"


REPOS=(
    "https://orbyta-isi-foundation@dev.azure.com/orbyta-isi-foundation/influenzanet/_git/api"

    "https://orbyta-isi-foundation@dev.azure.com/orbyta-isi-foundation/influenzanet/_git/user-management-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/study-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/messaging-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/logging-service"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/infra"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/api-gateway"
    "https://dev.azure.com/orbyta-isi-foundation/influenzanet/_git/participant-webapp"

    "https://orbyta-isi-foundation@dev.azure.com/orbyta-isi-foundation/influenzanet/_git/case-web-app-core"
    "https://orbyta-isi-foundation@dev.azure.com/orbyta-isi-foundation/influenzanet/_git/case-web-ui"
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
            git checkout main
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


cp ".env.development.local" "../../participant-webapp/"
echo "Copied .env.development.local to participant-webapp."

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

if [ -f "./script/import-db.sh" ]; then
    echo "Running import-db.sh..."
    chmod +x ./script/import-db.sh
    ./script/import-db.sh
else
    echo "Skipping database import, file not found."
    sleep 5
fi

echo "InfluenzaNet is ready to GO 🚀🚀🚀🚀🚀🚀"

cd ../participant-webapp
yarn install

echo "Starting frontend on new tab..."

GIT_BASH_PATH="C:\\Program Files\\Git\\git-bash.exe"

PARTICIPANT_WEBAPP_PATH=$(realpath ../participant-webapp)

"$GIT_BASH_PATH" -c "cd \"$PARTICIPANT_WEBAPP_PATH\" && yarn start" &

echo "/src/configs" >> .gitignore
echo "/src/configs/*" >> .gitignore
