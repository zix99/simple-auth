# Decoding a JWT Cookie

There are many ways to [decode a jwt cookie](https://jwt.io/).  This is an example
of using the [same-domain cookie](/access/cookie) in order to validate requests.

If request validation failed, they will be redirected to the auth-portal (*simple-auth*) to signin.

<a :href="`${$themeConfig.fileUrl}/docs/examples/traefik-subdomain/testapp`" target="_blank">View the full source code</a>

## Source

### package.json

<<< @/examples/traefik-subdomain/testapp/package.json

### index.js

<<< @/examples/traefik-subdomain/testapp/index.js
