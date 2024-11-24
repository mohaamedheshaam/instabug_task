#!/bin/bash
set -e

# Wait for MySQL
until nc -z db 3306; do
  echo "Waiting for MySQL to be ready..."
  sleep 2
done

# Setup database if not exists
bundle exec rails db:prepare

# Remove a potentially pre-existing server.pid for Rails.
rm -f /app/tmp/pids/server.pid

# Then exec the container's main process
exec "$@"