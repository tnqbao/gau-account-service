#!/bin/sh

# Get service type from first argument, default to "http"
SERVICE_TYPE=${1:-http}

echo "Starting service: $SERVICE_TYPE"

# Run migrations first for both services (common step)
echo "Running migrations..."
migrate -database "$PGPOOL_URL" -path migrations up
if [ $? -ne 0 ]; then
    echo "Migrations failed. Exiting."
    exit 1
fi
echo "Migrations completed successfully."

# Start the appropriate service
if [ "$SERVICE_TYPE" = "consumer" ]; then
    echo "Starting Consumer service..."
    ./consumer-service
else
    # Default to HTTP service
    echo "Starting HTTP API service..."
    ./http-service
fi