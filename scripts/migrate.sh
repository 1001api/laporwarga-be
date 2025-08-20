#!/bin/bash

# This script provides a simple interface for running database migrations.
# It ensures that everyone uses the same commands and paths.
# Exit immediately if a command exits with a non-zero status.
set -e

# Load environment variables from a .env file if it exists in the root.
if [ -f "$(dirname "$0")/../.env" ]; then
  while IFS='=' read -r key value; do
    if [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
      export "$key=$value"
    fi
  done < <(grep -v '^#' "$(dirname "$0")/../.env" | sed '/^\s*$/d')
fi

# Check if DATABASE_URL is set
if [ -z "${DATABASE_URL}" ]; then
  echo "Error: DATABASE_URL environment variable is not set."
  exit 1
fi

MIGRATION_PATH="internal/database/migrations"

# The first argument to the script determines the action
ACTION=$1

case "$ACTION" in
  up)
    echo "Applying all 'up' migrations..."
    migrate -database "${DATABASE_URL}" -path "${MIGRATION_PATH}" up
    ;;
  down)
    echo "Applying one 'down' migration..."
    migrate -database "${DATABASE_URL}" -path "${MIGRATION_PATH}" down 1
    ;;
  goto)
    VERSION=$2
    if [ -z "$VERSION" ]; then
      echo "Error: Please specify a version number for 'goto'."
      exit 1
    fi
    echo "Migrating to version ${VERSION}..."
    migrate -database "${DATABASE_URL}" -path "${MIGRATION_PATH}" goto "${VERSION}"
    ;;
  force)
    VERSION=$2
    if [ -z "$VERSION" ]; then
      echo "Error: Please specify a version number for 'force'."
      exit 1
    fi
    echo "Forcing migration to version ${VERSION}..."
    migrate -database "${DATABASE_URL}" -path "${MIGRATION_PATH}" force "${VERSION}"
    ;;
  drop)
    echo "Dropping all tables..."
    migrate -database "${DATABASE_URL}" -path "${MIGRATION_PATH}" drop
    ;;
  *)
    echo "Usage: $0 {up|down|drop|goto <version>}"
    exit 1
    ;;
esac

echo "Migration command finished."
