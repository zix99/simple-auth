# Simple Auth

*Simple-Auth* is a designed to be an easy way to manage users in a multi-usecase scenario.

**Common use-cases include:**

- Stand-alone user management, and credential validation via API
- Reverse proxy to downstream service blocked by login
- Per-request vouching (eg. for NGINX `auth_request` to act as a validator for login)
- Same-domain/subdomain login provider via cookie validation
- Various API implementations to authenticate a user (Simple, Three-way-Token)

**And it providers the common functionality:**

- "Simple" credentials (Username and password)
- OpenID Connect Login
- User email verification
- Forgot/lost password
- Login/access-attempt auditing

## Concepts

By default, simple-auth is split into three layers:

1. Login Providers
1. Authentication Providers (API)
1. Gateway/Login interface

## Login Providers

Login providers are how the user logs into simple-auth.  For example, they could be a "Simple" account (Username and password) or OpenID Connect (OAuth 2)

## Authentication Providers

Authentication providers are generally APIs where a downstream service can authenticate against *simple-auth* via API/redirect requests.

## Access

### Cookie

**Login-mode** is when *simple-auth* acts as a login-provider, and a downstream service validates the user via signed cookie or other API mechanism


### Gateway
**Gateway-mode** is when *simple-auth* sits between the user and what they're trying to access as a portal
