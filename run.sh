#!/bin/sh

echo "Launching Hitched API..."

# Load environment file.
if [ -f ./.env ]; then
    set -a
    . ./.env
    set +a
else
    echo "WARNING: No environment file found..."
fi

# Launch app
go run main.go
