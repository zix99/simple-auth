# Vouch API

The vouch API allows a down-stream user to check if the request is pre-authenticated (meaning the user has a session).  This will only work
from same-domain site.  It is intended to be used with [nginx auth_request](../cookbooks/nginx-auth-request), but may suite other use-cases.

## Enabling

```yaml
authenticators:
    vouch:
        enabled: true
```

## Making a request

From the same-domain (making sure cookies are passed) make a **GET** request to `/api/v1/auth/vouch`.  It will return 200 on success, otherwise 401

## See Also

- [nginx auth_request](/cookbooks/nginx-auth-request.md) - Using vouch with nginx
