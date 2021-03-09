# Simple Auth

<a :href="$themeConfig.repoUrl" target="_blank"><img src="https://img.shields.io/github/stars/zix99/simple-auth?style=social" title="Repo stars"></a>
<a :href="$themeConfig.repoUrl" target="_blank"><img src="https://img.shields.io/github/watchers/zix99/simple-auth?style=social" title="Repo Watchers"></a>

*Simple-Auth* is a designed to be an easy way to manage your site's users.  Unlike large complex solutions, it aims
to provide **simple login and user-management to a small or medium site**.  It doesn't try to replace global authentication providers
or enterprise user management (kerberos, active directory, etc...)

<a href="https://github.com/zix99/rare/releases" target="_blank"><img src="https://img.shields.io/github/v/release/zix99/simple-auth" alt="GitHub release (latest by date)"></a>
<a href="https://github.com/zix99/rare/releases" target="_blank"><img src="https://img.shields.io/github/downloads/zix99/simple-auth/total" alt="GitHub all releases"></a>
<a href="https://hub.docker.com/r/zix99/simple-auth" target="_blank"><img src="https://img.shields.io/docker/pulls/zix99/simple-auth" alt="Docker Pulls"></a>
<a href="https://hub.docker.com/r/zix99/simple-auth" target="_blank"><img src="https://img.shields.io/docker/image-size/zix99/simple-auth/latest" alt="Docker Image Size (latest)"></a>
<a href="https://github.com/zix99/rare/releases" target="_blank"><img src="" alt=""></a>
<a href="https://github.com/zix99/rare/releases" target="_blank"><img src="" alt=""></a>
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zix99/simple-auth)
![Coverage](./coverage.svg)

::: tip
Looking to get started? See [Quickstart](quickstart)
:::

![Simpleauth](./simpleauth.png)

**Features include:**

- Local user management (create account, login, [TOTP two-factor](/login/local.md#totp-2fa))
- [Credential validation via API](/authenticators/simple.md) and [OAuth2 / OIDC](/authenticators/oauth2.md)
- Reverse proxy to downstream service blocked by login ([gateway](/access/gateway.md))
- Per-request [vouching](/authenticators/vouch.md) (eg. for NGINX `auth_request` to act as a validator for login), to act as an authentication portal
- [Same-domain/subdomain login](/access/cookie.md) provider via cookie validation
- Various API implementations to [authenticate a user](/login) (Local username/password, third-party OIDC, etc)
- [OpenID Connect Login](/login/oidc.md) (OIDC) eg. Google Auth
- Optional [welcome email and email-verification](/email.md)
- Forgot/lost password
- Login/access-attempt auditing
- <a :href="`${$themeConfig.docsUrl}/api`">REST API</a> to all underlying functionality
- Mobile friendly
- White-label deployment using [customizations](/customization.md)

## Why Not...

There are plenty of other authentication providers out there.  You can always roll your own or use another solution like [Okta](https://www.okta.com/), [Gluu](https://www.gluu.org/) or [Keycloak](https://www.keycloak.org/).  While these services are perfectly fine (they're great, infact), *simple-auth* tries to be ***simple***.  Our [quickstart](quickstart.md) is incredibly short and the hosting modes allow **zero-to-fully setup in less than 5 minutes**.

Long story short, if you have the use case and time to look at another provider, please do! If you're looking for something simple and easy to get started with, *simple-auth* may be for you.

## Concepts

### Objects

At the root of the object representation sits the "Account". It is associated with
a unique email.  By itself, an account does not give access to login, it needs
an authentication object associated with it.

<mermaid>
graph TD
A[Account] --> AA{User Authenticates}
subgraph Authenticators
  AA --> B[Local Auth]
  B --> B2[TOTP/2FA]
  AA --> C[OAuth2]
  AA --> E[One-Time Auth]
end
A --> D[Audit Log]
</mermaid>

By default, simple-auth is split into three layers:

1. **Login Providers**: The mechanisms that allow creating an account object, and how a user logs in.  For example, a Local account or OIDC (OAuth2)
1. **Authenticators (API)**: How dowstream apps can authenticate with *simple-auth*. Usually via API or requests
1. **Access Layer**: How web applications can authenticate with simple-auth

You can find more information on all three layers on the left. Not all 3 layers are required for a functional setup.

# Next Steps

Try heading over to [Quickstart](quickstart) and giving it a try!
