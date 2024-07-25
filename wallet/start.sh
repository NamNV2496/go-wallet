#!/bin/sh

# Wait for Kafka to be available
while ! nc -z kafka 29092; do
  echo "Waiting for Kafka..."
  sleep 2
done

# Start the main application
exec "$@"
