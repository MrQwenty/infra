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
        if [ -f "docker-compose.yml" ]; then
            echo "Running Docker Compose in $REPO_NAME..."
            docker-compose up -d
        else
            echo "No docker-compose.yml found in $REPO_NAME, skipping..."
        fi
done

echo "Script executed successfully! All repositories have been processed."