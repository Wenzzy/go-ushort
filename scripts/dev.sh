#!/bin/bash

DEV_DOCKER_COMPOSE_BASE_FILE_NAME="./deploy/development/docker-compose"
DEV_DOCKER_ENV_BASE_FILE_NAME="./deploy/development/.env"
DEV_ENV_BASE_FILE_PATH="./configs/.env"
DEV_DB_DIAGRAM_PATH="./assets/db-diagram"
DOCKER_PROJECT_NAME="ushorter"

if [ -f ${DEV_ENV_BASE_FILE_PATH} ]; then
	# only on dev
	# shellcheck disable=SC1090
	source ${DEV_ENV_BASE_FILE_PATH}
fi

MIGRATIONS_PATH="./app/common/database/migrations"
MIGRATIONS_DB_URI="postgresql://${DB_USER:-dev_user}:${DB_PASS:-dev}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-$DOCKER_PROJECT_NAME}"

ENTRYPOINT_FILE="./app/cmd/app/main.go"
BUILT_FILE_PATH="./dist/"

usage() {
	cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") arg1 [arg2...]

Available options:

docker.all        Run app with all needed services in docker (build required)
docker.all-nb     Run app with all needed services in docker (NO-build)
docker.db         Run database in docker only
docker.db-metrics Run database and metrics (vm, grafana) in docker only
docker.stop       Stop project docker-compose
run               Run server
build             build app (to folder: "./dist")
install           install dependencies
test              Run tests of app
mg:r              Apply migrations
mg:rv             Rollback migrations
mg:c [name]       Create migration with [name]
docs              Fmt comments and generate swagger docs static files
db_diagram        Run database diagram generator (java required, database must be running)

EOF
	exit
}

run_logs() {
	trap 'sh run.sh docker.stop' INT
	#  docker compose -p ${DOCKER_PROJECT_NAME} logs -f app
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

generate_db_diagram() {
	java \
		-jar ./bin/schemaspy-6.2.3.jar \
		-t pgsql \
		-db ${DB_NAME:-$DOCKER_PROJECT_NAME} \
		-host ${DB_HOST:-localhost} \
		-u ${DB_USER:-dev_user} \
		-p ${DB_PASS:-dev} \
		-vizjs \
		-o ./tmp_schema \
		-imageformat svg \
		-dp ./bin/postgresql-42.5.4.jar \
		-s public
	mv -f ./tmp_schema/diagrams ${DEV_DB_DIAGRAM_PATH}
	mv -f ./tmp_schema/images ${DEV_DB_DIAGRAM_PATH}
	rm -rf ./tmp_schema
}

case $1 in
'docker.all')
	docker compose -f "${DEV_DOCKER_COMPOSE_BASE_FILE_NAME}.dev.yml" --env-file "${DEV_DOCKER_ENV_BASE_FILE_NAME}.dev" -p ${DOCKER_PROJECT_NAME} up -d -V --build --force-recreate
	run_logs
	;;
'docker.all-nb')
	docker compose -f "${DEV_DOCKER_COMPOSE_BASE_FILE_NAME}.dev.yml" --env-file "${DEV_DOCKER_ENV_BASE_FILE_NAME}.dev" -p ${DOCKER_PROJECT_NAME} up -d -V
	run_logs
	;;
'docker.db')
	docker compose -f "${DEV_DOCKER_COMPOSE_BASE_FILE_NAME}.db_only.yml" --env-file "${DEV_DOCKER_ENV_BASE_FILE_NAME}.dev" -p ${DOCKER_PROJECT_NAME} up -d --force-recreate
	;;
'docker.db-metrics')
	docker compose -f "${DEV_DOCKER_COMPOSE_BASE_FILE_NAME}.db_metrics_only.yml" --env-file "${DEV_DOCKER_ENV_BASE_FILE_NAME}.dev" -p ${DOCKER_PROJECT_NAME} up -d --force-recreate
	;;
'docker.stop')
	docker compose -p ${DOCKER_PROJECT_NAME} down
	;;
'run')
	go run ${ENTRYPOINT_FILE}
	;;
'build')
	go build -o ${BUILT_FILE_PATH} ${ENTRYPOINT_FILE}
	;;
'install')
	go install ${ENTRYPOINT_FILE}
	;;
'test')
	run_tests
	;;
'mg:r')
	goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} up
	;;
'mg:rv')
	goose -dir ${MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} down
	;;
'mg:c')
	goose -dir ${MIGRATIONS_PATH} create $2 sql
	;;
'docs')
	swag fmt -g ${ENTRYPOINT_FILE}
	swag init --parseInternal --parseDepth 2 -g ${ENTRYPOINT_FILE}
	;;
'db_diagram')
	generate_db_diagram
	;;
*)
	usage
	;;
esac
