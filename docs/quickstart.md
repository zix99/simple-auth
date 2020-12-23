# Quick-Start

[[toc]]

## Stand-Alone Authentication Portal

The easiest way to get started is to run *simple-auth* as a stand-alone portal, accessed by API, to manage users.

By default, only one config value is required to be set: `web.login.cookie.jwt.signingkey`  This key must be unique, cryptographically secure, and very secret, to guarantee a user's session can't be hijacked.

Once the application is started, it will create a local sqlite3 `simpleauth.db` file to store the users.  You can optionally change this to a different database. See [database providers](/database)

::: danger
The `signingkey` must be kept secret at all times. This is how a user can login, and
*simple-auth* knows who they are.  If you need to share the key to validate the JWT,
I recommend using public-private [key pair strategy](config.md#signing-key-pair) (RS256, RS512)
:::

::: warning
The app will start without a `signingkey` but the user won't be able to login
:::

### Docker

::: tip
By default, simple-auth in docker will put your database in `/var/lib/simple-auth`. Make sure to create a volume so you don't lose your data
:::

```sh
docker run -it --rm -e SA_WEB_LOGIN_COOKIE_JWT_SIGNINGKEY=REPLACE_ME zix99/simple-auth
```

Or, if you prefer *docker-compose*...

<<< @/examples/simple/docker-compose.yml

### Binary

Download the binary from the releases page, and run with:

```sh
./simple-auth-server --web-login-cookie-jwt-signingkey=REPLACE_ME
```

## Next Steps

Check out [more config options](config), [customize the UI](customization), or check out some [cookbooks](cookbooks/)

You can also check out the <a :href="`${ $themeConfig.repoUrl }/docs/examples`">examples in the repository</a>