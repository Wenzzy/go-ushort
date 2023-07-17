#!/bin/bash


DEV_COMPOSE_BASE_FILE_NAME="./deploy/development/docker-compose"
DEV_ENV_BASE_FILE_NAME="./deploy/development/.env"
PROJECT_NAME="ushorter"

MIGRATIONS_PATH="common/database/migrations"
MIGRATIONS_DB_URI="postgresql://dev_user:dev@localhost:5432/ushorter"

ENTRYPOINT_FILE="./app/cmd/app/main.go"

run_logs() {
  trap 'sh run.sh stop' INT
#  docker compose -p ${PROJECT_NAME} logs -f app
  docker logs -f app
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
'docker.stop')
  docker compose -p ${PROJECT_NAME} down;;
'run')
  go run ${ENTRYPOINT_FILE};;
'mg:r')
  goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} up;;
'mg:rv')
  goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} down;;
'docs')
  swag fmt
  swag init
  ;;
*)
  echo "please run with not empty attributes (docker.all, docker.all-nb, docker.db, docker.stop, run, mg:r, mg:rv, docs)";;
esac
