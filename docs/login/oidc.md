# OpenID Connect (OIDC)

OpenID Connect Login is a subset of OAuth 2.0 spec that allows users to easily login to many third-party providers.

Once setup, there will be a button presented below the login form on homepage to "Connect to *OIDC Name*"

## Configuring Google OIDC

You will need to obstain an OIDC id and secret from [google OIDC page](https://developers.google.com/identity/protocols/oauth2/openid-connect)

```yaml
web:
  login:
    oidc:
    - id: google     # Unique ID (stored in DB as provider name)
      name: 'Google' # Presented in the ui
      icon: 'google' # Font-awesome icon
      clientid: '<CLIENT ID>'
      clientsecret: '<CLIENT SECRET>'
      authurl: 'https://accounts.google.com/o/oauth2/v2/auth'
      tokenurl: 'https://oauth2.googleapis.com/token'
```
