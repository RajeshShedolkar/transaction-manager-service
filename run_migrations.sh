#!/bin/bash

set -e

MIGRATIONS_DIR="$(pwd)/migrations"
DB_DSN="postgres://transaction_service:pass123@host.docker.internal:5432/transactiondb?sslmode=disable"

echo "🚀 Running DB migrations..."
echo "Using DSN: $DB_DSN"

docker run --rm \
  -v "$MIGRATIONS_DIR:/migrations" \
  migrate/migrate \
  -path /migrations \
  -database "$DB_DSN" \
  up

echo "✅ Migrations completed successfully."
