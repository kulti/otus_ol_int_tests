#!/bin/bash

set -e

cd "$(dirname "$0")"

env=$1
cmd=$2
shift 2 || true

project_name=${PROJECT_NAME:-${env}}

compose="docker-compose -p ${project_name} -f docker-compose.yml -f docker-compose.${env}.yml"

case "${cmd}" in
up)
    services=$*
    eval "${compose} run --rm game_server_createdb"
    eval "${compose} run --rm game_server_migratedb"
    eval "${compose} run --rm user_stats_createdb"
    eval "${compose} run --rm user_stats_migratedb"

    eval "${compose} up -d ${services}"
    ;;
down|ps|logs|exec|start|stop|restart)
    args=$*
    eval "${compose} ${cmd} ${args}"
    ;;
recreate)
    services=$*
    eval "${compose} up -d --no-deps ${services}"
    ;;
run)
    service=$1
    eval "${compose} run --rm ${service}"
    ;;
*)
cat << EOF
Usage ./env.sh <environemnt> <command> [args...]

Environments:
    dev  - used for local development.
           Examples:
           ./env.sh dev up
           ./env.sh dev ps
           ./env.sh dev logs
           ./env.sh dev logs game_server
    test - used for running integration tests on CI.

Commands:
    up, down, recreate
    run, exec
    start, stop, restart
    ps, logs
EOF
    exit 1
    ;;
esac
