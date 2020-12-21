# Simple Auth API

::: warning
While the simple API, is.. well.. simple, it only support the **local provider**.
:::

The Simple API providers a way to authenticate a username, password, and optional 2FA code, against an endpoint.

## Enabling

```yaml
authenticators:
    simple:
        enabled: true # set to true
        sharedsecret: null # an optional shared-secret (but you should use it!)
```

## Using

**Endpoint:** `POST /api/v1/auth/simple`

**Body:**
```json
{
  "username": "chris",
  "password": "super-secret",
  "totp": "123456" // Optional
}
```

**Security:**
The endpoint is secured using an optional shared-secret. When provided, pass it to the endpoint via a header:

```http
Authorization: Bearer <shared-secret>
```

**Returns:**
- `200` with a simple payload `{"id": "uuid"}` on success
- `403` on failure (bad credentials)
- `401` when you're not allowed to hit this endpoint