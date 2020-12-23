# Database Drivers

[[toc]]

By default, *simple-auth* stores all credentials as a *sqlite3* database named `simpleauth.db` in the working directory of the app.
In docker the docker file, that may be in a slightly different spot (See docker docs)

In some cases, you may want to change that.  For instance, postgres often performs better than sqlite.  Or, you might want
your persistence in a different spot.

## Configuring Simple-Auth

This is the applicable section of the configuration file (with its defaults)

```yaml
db:
    driver: "sqlite3"     # Storage driver: "sqlite3", "postgres", "mysql"
    url: "simpleauth.db"  # Storage connection URL. See http://gorm.io/docs/connecting_to_the_database.html
```

In order to use different persistence, you can chose a different driver and url.

## Drivers

### Sqlite

This is the default.  To sepecify a different path, change the `db.url` config.

#### Docker

If using docker, make sure to mount a persistent volume, otherwise you may lose your users on container restart.

*simple-auth* will put all its persistent storage here: `/var/lib/simple-auth`

<<< @/examples/simple/docker-compose.yml

### Sqlite in-memory

::: warning
You will lose all data upon application termination
:::

It can sometimes be useful to have a temporary in-memory database (mainly for testing or integration testing). For this
you can use the `sqlite3` driver, with the url `file::memory:?cache=shared` like so:

```yaml
db:
  url: "file::memory:?cache=shared"
```

### Postgres

This is an example setup using a docker-compose file to specify the config for a postgres instance.

If you use a standalone applicable, you only need to set `SA_DB_DRIVER` and `SA_DB_URL`

<<< @/examples/postgres/docker-compose.yml

### Mysql

::: warning
Make sure to set `?charset=utf8&parseTime=True&loc=Local` on the connection, otherwise you will receive errors
:::

<<< @/examples/mysql/docker-compose.yml