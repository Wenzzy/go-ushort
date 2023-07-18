# go-ushorter
My pet project that shortens long links ;)

## Todo
- [x] JWT Authentication (access,refresh)
- [x] Create/update link by user
- [ ] Configure JSON-only logger on IsProduction mode
- [x] Create script (alias) for migration creation
- [ ] Write tests


For simple run scripts - make alias `alias gr="sh run.sh"`


### General scripts

```shell
gr docker.all
# Run app with all needed services in docker (build required)
```

```shell
gr docker.all-nb
# Run app with all needed services in docker (NO-build)
```

```shell
gr docker.db
# Run database in docker only
```

```shell
gr docker.stop
# Stop project docker-compose
```

```shell
gr docs
# Fmt comments and generate swagger docs static files
```



### Database scripts

```shell
gr mg:r
# apply migrations
```


```shell
gr mg:rv
# rollback migrations
```

```shell
gr mg:c [name]
# create new migration
```

