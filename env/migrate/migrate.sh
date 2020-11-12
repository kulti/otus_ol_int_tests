#!/bin/sh

set -xue

migrate -path /app/migrations/${MIGRATION_ID} -database ${DB_URL} $*
