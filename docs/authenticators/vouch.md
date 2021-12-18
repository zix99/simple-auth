# Vouch API

The vouch API allows a down-stream user to check if the request is pre-authenticated (meaning the user has a session).  This will only work
from same-domain site.  It is intended to be used with [nginx auth_request](../cookbooks/nginx-auth-request) or [traefik forwardauth](../cookbooks/traefik.md),
but may suite other use-cases.

## Enabling

```yaml
authenticators:
    vouch:
        enabled: true
```

## Making a request

From the same-domain (making sure cookies are passed) make a **GET** request to `/api/v1/auth/vouch`.  It will return 200 on success, otherwise 401

## Forward Auth

By default, vouch returns a `401` if authentication doesn't exist. Enabling *forward* will instead return a temporary redirect via a `307` to the authentication
portal.  It checks for a `continue=` query param or the `X-Forwarded` headers traefik provides to send the user back to their source.  These continue
urls must be allow-listed via `web.login.settings.allowedcontinueurls`.

Similar to the vouch, this is a simple **GET** request to `/api/v1/auth/vouch?forward=1`

## See Also

- [nginx auth_request](/cookbooks/nginx-auth-request.md) - Using vouch with nginx
- [traefik forwardauth](/cookbooks/traefik.md) - Using vouch with traefik
