# Same-domain Cookie

One simple strategy for having *simple-auth* manage your authentication is to put *simple-auth* on the same
domain on a path or subdomain, and validating the JWT cookie on your app.

::: tip
If you're validating the JWT, you'll need to share the `signingkey` between *simple-auth* and your app. For
more security, please consider using an [RSA Key-Pair](/cookbooks/signingkey-pair).
:::

## Setting up for same-domain

The most important item for setting-up for same-domain is making sure that the cookie is on the correct
top-level domain.  For this example, we'll have our content on `example.com` and *simple-auth* will
be on `auth.example.com`

<mermaid>
graph LR
A{User} --> B[Nginx]
B -- auth.example.com --> C[Simple Auth]
B -- example.com --> D[Backend]
C -.-> D
</mermaid>

```yaml
web:
    baseurl: https://auth.example.com
    # Login section for setting different ways a user is allowed to login
    login:
        settings:
            routeonlogin: https://example.com
        cookie:
            jwt:
                signingkey: "" # INSERT YOUR KEY
            path: /                # The path the session will be stored at (mainly useful if simple-auth is at a sub-path of root)
            domain: example.com
```

## Cookie

::: tip
The cookie name can be set by `web.login.cookie.name`, and by default it is `auth`.
:::

::: danger
You must validate the cookie using a JWT library and the `signingkey`.  If you don't, it's possible for anyone to forge the cookie!
:::

Once logged-in, the cookie will be set according to the config in `web.login.cooke` as a JWT.  This cookie is also used to login and display the *simple-auth* management UI.


### Validating the JWT

Because anyone can form a JWT (not encrypted, just signed), you need to make sure to check the signature
in your application before processing any data.

Since there are so many languages, I'll refer you to [jwt.io](https://jwt.io/) which has numerous implementations and
examples.  You can also check out [RFC7519](https://tools.ietf.org/html/rfc7519)

## See Also

* [jwt.io](https://jwt.io/)
* [Same Domain Cookie Auth using Traefik](/cookbooks/traefik.md)
