# Gateway

[[toc]]

Gateway is a simple way for *simple-auth* to **act as an authentication gateway and reverse-proxy to a service**.  This is great for low-volume sites, private setups, or testing, but you probably don't want to use *simple-auth* as a reverse-proxy in high-volume situations (for many reasons).

::: tip
The gateway can target many backend servers, **however**, it can't route to different
servers depending on host or path.  If you have a more complex use-case, please
look at [nginx auth request](/cookbooks/nginx-auth-request).
:::

::: warning
If *simple-auth* acts as a reverse proxy to your site, that depends on you securing your site in a different way (firewall rules, docker network, etc).
:::

<mermaid>
graph LR
A{User} -- Web Request --> B[Simple-Auth]
B -- Proxies --> C[Backend]
</mermaid>

## Config

```yaml
web:
  gateway:
    enabled: true
    basicauth: false       # If true, will allow basic auth for local auth to pass-through the gateway
    logoutpath: "/logout"  # Special path that will act as "logout" (clear session).  Shouldn't conflict with any downstream URLs
    targets:               # One or more downstream servers that SA will proxy to
      - example.com
    host: example.com      # Override the host header
    rewrite: null          # Rewrite URLs upon proxying eg "/old"->"/new" or "/api/*"->"/$1"
    headers: null          # Write additional headers (excluding host header)
    nocache: true          # If true, will attempt to disable caching to gateway target
```

For example:

<<< @/examples/sh/gateway.sh

::: tip
If you aren't seeing the downstream page you expect, you might need to set the `host` header
via the `web.gateway.host` config.
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

## Basic Auth

Basic auth is available so that making API calls or requests via non-browser on the gateway is possible.  It is disabled by default, but can be enabled by setting `web.gateway.basicauth` to `true`.

::: warning TOTP/2FA
If the account has TOTP enabled, basic-auth via gateway will not be supported.
:::

This allows you to make the following request to the downstream service with credentials:

```bash
curl http://apiuser:apipass@example.com
```

The normal headers will be passed to the downstream service **except** the `Authorization` header will be removed from the request.

