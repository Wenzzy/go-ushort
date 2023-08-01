<p align="center">
 <img src='assets/logo.jpg' width='500'>
</p>

[![CI-CD pipeline](https://github.com/WenzzyX/go-ushort/actions/workflows/ci-cd.production.yml/badge.svg)](https://github.com/WenzzyX/go-ushort/actions/workflows/ci-cd.production.yml)
&nbsp;\
My pet project that shortens long links ;)\
&nbsp;\
Front-end: [go-ushort](https://github.com/WenzzyX/frontend-ushort)\
Design: [figma](https://www.figma.com/community/file/1267929060708092881)
&nbsp;\
&nbsp;\
_`https://very-long-site-subdomain.long-domain-ffff.com/my-best-blog`_ \
-> `https://ushort.us/7R`

## Todo

- [x] JWT Authentication (access,refresh)
- [x] Create/update link by user
- [x] Configure JSON-only logger on IsProduction mode
- [x] Create script (alias) for migration creation
- [x] Add pagination for links
- [x] Add a collection of metrics and a dashboard to monitor them (prom-app -> victoriametrics -> grafana)
- [x] Add LICENSE
- [x] Write tests
- [x] Write the rules for making a contribution
- [x] Change db diagram svg generation (replace images)
- [x] Configure CI/CD
- [x] Configure dev_full docker-compose
- [x] Configure deploy (render.com)
- [x] Link domain
- [ ] Configure git-crypt

For simple run scripts - make alias `alias gr="sh run.sh"`

### General

```shell
gr
# or `sh run.sh`
# Get list of scripts and description for each script
```

```shell
docker build . \
	--platform=linux/amd64 \
	-t go-ushort \
	--build-arg NEXT_PUBLIC_USHORT_DOMAIN="ushort.us"
# build for amd64

docker build . -t go-ushort
# build docker image
```

```shell
docker rm go-ushort && \
docker run -it -p 5005:8000 \
	--name go-ushort \
	-e ALLOWED_ORIGINS="https://go-ushort.vercel.app" \
	-e DB_HOST="localhost" \
	-e DB_NAME="ushorter" \
	-e DB_PASS="dev" \
	-e DB_PORT="5432" \
	-e DB_USER="dev_user" \
	-e DOMAIN="ushort.us" \
	-e IS_DEBUG="false" \
	-e IS_ENABLE_PROM="false" \
	-e IS_PRODUCTION="true" \
	-e JWT_ACCESS_EXP_TIME="1m" \
	-e JWT_ACCESS_SECRET="ocSbpF5qQjBbutPR85g7VHfQn1v7dGYO0IVEoH9xq2hmWDa6bVxX8NWk6OcpdEZN" \
	-e JWT_REFRESH_EXP_TIME="30d" \
	-e JWT_REFRESH_SECRET="mwgqOZsFf8hWNdOtbKQQLGPhwWXQQQ0hHOKZvypj82uJuENwjNPqXLBMdKRYsqBq" \
	-e MIGRATIONS_PATH="./migraions" \
	go-ushort \
	&& docker logs -f go-ushort
# run docker container
```

### Environment

#### Server config

| param                     | type      | required | default   | description                                                      |
| ------------------------- | --------- | -------- | --------- | ---------------------------------------------------------------- |
| `JWT_ACCESS_SECRET`       | `string`  | `yes`    | `-`       | Secret for generating accessToken                                |
| `JWT_ACCESS_EXP_TIME`     | `string`  | `yes`    | `-`       | life duration of accessToken (ex.: "20s", "2d")                  |
| `JWT_REFRESH_SECRET`      | `string`  | `yes`    | `-`       | Secret for generating refreshToken                               |
| `JWT_REFRESH_EXP_TIME`    | `string`  | `yes`    | `-`       | life duration of accessToken (ex.: "20s", "2d")                  |
| `IS_PRODUCTION`           | `boolean` | `no`     | `true`    | Run-mode is production? (ex.: "true")                            |
| `IS_DEBUG`                | `boolean` | `no`     | `false`   | Print sensitive info and prettify log messages? (ex.: "true")    |
| `IS_ENABLE_PROM`          | `boolean` | `no`     | `false`   | Enable prometheus? (ex.: "false")                                |
| `DOMAIN`                  | `string`  | `yes`    | `-`       | Domain for setting cookies (ex.: "localhost")                    |
| `ALLOWED_HOSTS`           | `string`  | `no`     | `0.0.0.0` | Hosts who can send requst to server (ex.: "0.0.0.0,192.168.1.1") |
| `ALLOWED_ORIGINS`         | `string`  | `no`     | `*`       | CORS - origin (ex.: "https://ushort.us,http://localhost:3000")   |
| `SERVER_HOST`             | `string`  | `no`     | `0.0.0.0` | Host, where server will run (ex.: "0.0.0.0")                     |
| `LIMIT_COUNT_PER_REQUEST` | `int`     | `no`     | `10`      | _Temporarly not using_                                           |

&nbsp;\
&nbsp;

#### Database config

| param             | type      | required | default                            | description                                                      |
| ----------------- | --------- | -------- | ---------------------------------- | ---------------------------------------------------------------- |
| `DB_NAME`         | `string`  | `yes`    | `-`                                | DB name                                                          |
| `DB_USER`         | `string`  | `yes`    | `-`                                | DB user                                                          |
| `DB_PASS`         | `string`  | `yes`    | `-`                                | DB password                                                      |
| `DB_HOST`         | `string`  | `no`     | `localhost`                        | DB host (ex.: "localhost")                                       |
| `DB_PORT`         | `int`     | `no`     | `5432`                             | DB port (ex.: "5432")                                            |
| `DB_LOG_MODE`     | `boolean` | `no`     | `false`                            | Output SQL and other query information? (ex.: "true")            |
| `DB_SSL_MODE`     | `string`  | `no`     | `false`                            | Use SSL mode (ex.: "disable", "enable")                          |
| `MIGRATIONS_PATH` | `string`  | `no`     | `./app/common/database/migrations` | Migrations folder path (ex.: "./app/common/database/migrations") |

### Database diagram

![DB-diagram image](assets/db-diagram/summary/relationships.real.large.svg)
