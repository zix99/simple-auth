# Simple-Auth Gateway

[[toc]]

Gateway is a simple way for *simple-auth* to **act as an authentication gateway and reverse-proxy to a service**.  This is great for low-volume sites, private setups, or testing, but you probably don't want to use *simple-auth* as a reverse-proxy in high-volume situations (for many reasons).

::: tip
The gateway can target many backend servers, **however**, it can't route to different
servers depending on host or path.  If you have a more complex use-case, please
look at [nginx auth request](nginx-auth-request).
:::

::: warning
If *simple-auth* acts as a reverse proxy to your site, that depends on you securing your site in a different way (firewall rules, etc).
:::

## Headers

### Simple-auth Added

By default, *simple-auth* will added the following headers:

* `X-SA-Account` will contain the UUID of the logged-in user's account

If the `host` config setting is provider, it will also added the `Host` header.

### Additional

By adding values to the `web.gateway.headers` map, you can add additional headers
to the proxied request.

**NOTE:** This cannot be used to override the host header.

## URL Rewriting

Since *simple-auth* is acting as a gateway, you may want it to rewrite some URLs. You
can do this via the config `web.gateway.rewrite`.

For example:
```yaml
web:
  gateway:
    enable: true
    rewrite:
      "/old": "/new"
      "/api/*": "/$1"
```

# Examples

## Config example

```yaml
web:
  gateway:
    enabled: true
    logoutpath: "/logout"  # Special path that will act as "logout" (clear session).  Shouldn't conflict with any downstream URLs
    targets:
      - example.com
    host: example.com
    nocache: true          # If true, will attempt to disable caching to gateway target
```

## docker-compose example

<<< @/../examples/gateway/docker-compose.yml

