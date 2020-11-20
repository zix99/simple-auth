# OpenID Connect (OIDC)

OpenID Connect Login is a subset of OAuth 2.0 spec that allows users to easily login to many third-party providers.

## Configuring Google OIDC

```yaml
web:
  login:
    oidc:
    - id: google
      name: 'Google'
      icon: 'google'
      clientid: '<CLIENT ID>'
      clientsecret: '<CLIENT SECRET>'
      authurl: 'https://accounts.google.com/o/oauth2/v2/auth'
      tokenurl: 'https://oauth2.googleapis.com/token'
```