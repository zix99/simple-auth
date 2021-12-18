# Traefik Simple-Auth

[Traefik](https://traefik.io/) has become a popular choice for load-balancing docker applications.

::: danger
This is an example config, and doesn't have SSL enabled by default. [Traefik supports SSL](https://doc.traefik.io/traefik/https/overview/).
Make sure to enable it so that username and password are encrypted in transit!
:::

In both cases, we use a same-domain cookie sharing technique, described [here](/access/cookie)

## Forward Auth

This strategy is similar to [nginx auth_request](nginx-auth-request.md), where traefik will forward the request
to *simple-auth*'s `vouch` endpoint to see if a user has a session (in this case, stored in a cookie).  Unlike
nginx's `auth_request`, the user should be forwarded to *simple-auth* by the `vouch` endpoint if there is no valid session.

This can also be used via traefik to forward a header, eg the user id, to the downstream service.

<mermaid>
graph LR
A[Traefik]
A -- ForwardAuth--> C[Simple-auth]
A -- example.com --> B[App]
</mermaid>

### docker-compose

<<< @/examples/traefik-forwardauth/docker-compose.yml

## Same-Domain Cookie in App

This strategy uses [same-domain cookie auth](/access/cookie.md) to authenticate the user inside the app, with no special load-balancer setup.

Here, we're using traefik to have both *simple-auth* and a *testapp* (validates the token in the cookie).  The test-app will forward to `auth.${DOMAIN}` if it doesn't detect an `auth` token.

<mermaid>
graph LR
A[Traefik]
A -- example.com --> B[Test App]
A -- auth.example.com --> C[Simple-auth]
</mermaid>

### docker-compose

<<< @/examples/traefik-subdomain/docker-compose.yml


### Test App

::: tip
You can find more information about the testapp [here](/cookbooks/decodejwt)
:::

This is a very simple nodejs app that will validate your `auth` cookie, or redirect
you to the authentication portal if it fails.

You can see the full app at our <a :href="`${$themeConfig.fileUrl}/docs/examples/traefik-subdomain/testapp`" target="_blank">repository</a>

<<< @/examples/traefik-subdomain/testapp/index.js
