#!/bin/bash


DEV_COMPOSE_BASE_FILE_NAME="./deploy/development/docker-compose"
DEV_ENV_BASE_FILE_NAME="./deploy/development/.env"
PROJECT_NAME="ushorter"

MIGRATIONS_PATH="./app/common/database/migrations"
MIGRATIONS_DB_URI="postgresql://dev_user:dev@localhost:5432/ushorter"

TESTS_DB_URI="postgresql://dev_user:dev@localhost:5432/ushorter"

ENTRYPOINT_FILE="./app/cmd/app/main.go"
BUILT_FILE_PATH="./dist/"
TESTS_PATH="./tests/"

run_logs() {
  trap 'sh run.sh docker.stop' INT
#  docker compose -p ${PROJECT_NAME} logs -f app
  docker logs -f app
}

run_tests() {
  docker run \
        --name postgres-tests \
        --rm \
        -e POSTGRES_DB=tests \
        -e POSTGRES_USER=tests \
        -e POSTGRES_PASSWORD=tests \
        -p 6543:5432 \
        -v postgres-tests-data:/var/lib/postgresql/data \
        -d postgres:15.2-alpine
    go test ./...
    docker stop postgres-tests
    docker volume rm postgres-tests-data
}

usage() {
  cat << EOF
Usage: $(basename "${BASH_SOURCE[0]}") arg1 [arg2...]

Available options:

docker.all        Run app with all needed services in docker (build required)
docker.all-nb     Run app with all needed services in docker (NO-build)
docker.db         Run database in docker only
docker.db-metrics Run database and metrics (vm, grafana) in docker only
docker.stop       Stop project docker-compose
run               Run server
build             build app (to folder: "./dist")
test              Run tests of app
mg:r              Apply migrations
mg:rv             Rollback migrations
mg:c [name]       Create migration with [name]
docs              Fmt comments and generate swagger docs static files

EOF
  exit
}


case $1 in
'docker.all')
  docker compose -f "${DEV_COMPOSE_BASE_FILE_NAME}.dev.yml" --env-file "${DEV_ENV_BASE_FILE_NAME}.dev" -p ${PROJECT_NAME} up -d -V --build --force-recreate
  run_logs
  ;;
'docker.all-nb')
  docker compose -f "${DEV_COMPOSE_BASE_FILE_NAME}.dev.yml" --env-file "${DEV_ENV_BASE_FILE_NAME}.dev" -p ${PROJECT_NAME} up -d -V
  run_logs
  ;;
'docker.db')
  docker compose -f "${DEV_COMPOSE_BASE_FILE_NAME}.db_only.yml" --env-file "${DEV_ENV_BASE_FILE_NAME}.dev" -p ${PROJECT_NAME} up -d --force-recreate;;
'docker.db-metrics')
  docker compose -f "${DEV_COMPOSE_BASE_FILE_NAME}.db_metrics_only.yml" --env-file "${DEV_ENV_BASE_FILE_NAME}.dev" -p ${PROJECT_NAME} up -d --force-recreate;;
'docker.stop')
  docker compose -p ${PROJECT_NAME} down;;
'run')
  go run ${ENTRYPOINT_FILE};;
'build')
  go build -o ${BUILT_FILE_PATH} ${ENTRYPOINT_FILE} ;;
'test')
  run_tests;;
'mg:r')
  goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} up;;
'mg:rv')
  goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} down;;
'mg:c')
  goose -dir ${MIGRATIONS_PATH} create $2 sql;;
'docs')
  swag fmt -g ./app/cmd/app/main.go
  swag init --parseInternal --parseDepth 2 -g ./app/cmd/app/main.go
  ;;
*)
  usage;;
esac
