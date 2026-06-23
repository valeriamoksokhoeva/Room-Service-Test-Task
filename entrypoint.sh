#!/bin/sh
set -e

echo "Generating dbconfig.yml..."

cat > dbconfig.yml <<EOF
production:
    dialect: postgres
    datasource: "host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=disable"
    dir: migrations
    table: migrations
EOF

echo "Waiting for database..."

until pg_isready -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT}" -U "${POSTGRES_USER}"; do
  echo "Database is unavailable - sleeping"
  sleep 1
done

echo "Database is up - running migrations"

sql-migrate up -env=production

echo "Migrations applied - starting app"

exec ./app