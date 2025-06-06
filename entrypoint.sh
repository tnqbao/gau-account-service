#!/bin/sh
# Run migrations for the application
echo "Running migrations..."
migrate -database "$PGPOOL_URL" -path migrations up
if [ $? -ne 0 ]; then
    echo "Migrations failed. Exiting."
    exit 1
fi
echo "Migrations completed successfully."

# Start the application
echo "Starting API..."
if [ -f "gau-account-service.bin" ]; then
    ./gau-account-service.bin
else
    echo "Running main.go..."
    go run main.go
fi