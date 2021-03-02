# Quick-Start

[[toc]]

## Stand-Alone Authentication Portal

The easiest way to get started is to run *simple-auth* as a stand-alone portal, accessed by API, to manage users.

By default, *simple-auth* will run without any configuration. However, if you want the user to be able to login via cookie, you must set: `web.login.cookie.jwt.signingkey`  This key must be unique, cryptographically secure, and very secret, to guarantee a user's session can't be hijacked.

Once the application is started, it will create a local sqlite3 `simpleauth.db` file to store the users.  You can optionally change this to a different database. See [database providers](/database)

::: tip
A simple way to generate a secure password is with `openssl rand -base64 14`
:::

::: danger
The `signingkey` must be kept secret at all times. This is how a user can login, and
*simple-auth* knows who they are.  If you need to share the key to validate the JWT,
I recommend using public-private [key pair strategy](config.md#signing-key-pair) (RS256, RS512)
:::

### Docker

::: warning
By default, simple-auth in docker will put your database in `/var/lib/simple-auth`. Make sure to create a volume so you don't lose your data on container restart!
:::

```sh
docker run -it --rm -e SA_WEB_LOGIN_COOKIE_JWT_SIGNINGKEY=REPLACE_ME -p 80:80 zix99/simple-auth
```

Or, if you prefer *docker-compose*...

<<< @/examples/simple/docker-compose.yml

### Binary

Download the binary from the <a :href="`${$themeConfig.repoUrl}/releases`" target="_blank">releases page</a>, and run with:

```sh
./simple-auth-server --web-login-cookie-jwt-signingkey=REPLACE_ME
```

All environment variables can be replaced with CLI counterparts. For more information see [Config](/config.md)

## Simple Gateway

After running it stand-alone, the simplest way to start using authentication is the [Gateway](/access/gateway.md) functionality.  This proxies
request to a downstream HTTP service through *simple-auth* when the user is authenticated. To use this functionality you need 3 new environment variables:

```bash
SA_WEB_LOGIN_SETTINGS_ROUTEONLOGIN=/
SA_WEB_GATEWAY_ENABLED=true
SA_WEB_GATEWAY_TARGETS=http://downstream-target
```

For example:

<<< @/examples/sh/gateway.sh

Read more in the [Gateway](/access/gateway.md) docs

## TLS

If *simple-auth* is exposed to the public internet, you should use TLS encryption.  If you use a proxy, it might provide TLS for you.  If not, *simple-auth* has the ability to issue a valid certificate via [Let's Encrypt](https://letsencrypt.org/)

To enable, simply set `web.tls.enabled` to `true`.

For more information, see [TLS Cookbook](/cookbooks/letsencrypt.md)

## Next Steps

Check out [more config options](config), [customize the UI](customization), or check out some [cookbooks](cookbooks/)

You can also check out the <a :href="`${ $themeConfig.fileUrl }/docs/examples`">examples in the repository</a>