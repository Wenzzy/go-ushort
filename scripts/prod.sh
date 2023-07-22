#!/bin/bash

_MIGRATIONS_PATH="${MIGRATIONS_PATH:-migrations}"
MIGRATIONS_DB_URI="postgresql://${DB_USER:-dev_user}:${DB_PASS:-dev}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-$DOCKER_PROJECT_NAME}"

usage() {
	cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") arg1 [arg2...]

Available options:

mg:r              Apply migrations
mg:rv             Rollback migrations

EOF
	exit
}

case $1 in
'mg:r')
	goose -dir ${_MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} up
	;;
'mg:rv')
	goose -dir ${_MIGRATIONS_PATH} postgres ${MIGRATIONS_DB_URI} down
	;;
*)
	usage
	;;
esac
