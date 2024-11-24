#!/bin/bash
set -e

# Wait for MySQL
until nc -z db 3306; do
  echo "Waiting for MySQL to be ready..."
  sleep 2
done

# Wait for elasticsearch
until nc -z elasticsearch 9200; do
  echo "Waiting for Elasticsearch to be ready..."
  sleep 2
done

# Wait for RabbitMQ
until nc -z rabbitmq 5672; do
  echo "Waiting for RabbitMQ to be ready..."
  sleep 2
done

# Execute main command
exec "$@"