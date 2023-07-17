### Database actions

```shell
goose -dir common/database/migrations postgres "postgresql://dev_user:dev@localhost:5432/ushorter" up
# apply migrations
```


```shell
goose -dir common/database/migrations postgres "postgresql://dev_user:dev@localhost:5432/ushorter" down
# rollback migrations
```

```shell
goose -dir common/database/migrations create migration_name sql 
# create new migration
```

> for simple run scripts - make alias `alias gr="sh run.sh"`