#!/bin/sh

set -xue

DBNAME=$1

wait_for_db()
{
    for i in $(seq 1 30); do
        echo "SELECT 1" | psql -h postgres -U postgres && return
        sleep 1
    done

    exit 1
}

wait_for_db

echo "SELECT 'CREATE DATABASE ${DBNAME}' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DBNAME}')\gexec" | psql -h postgres -U postgres -v ON_ERROR_STOP=1
