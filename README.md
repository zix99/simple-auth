# simple-auth

Simple-auth is a lightweight, whitelabeled, authentication solution.  It allows users to sign-up
with a UI, offering various security measures.  Then, via various authenticators, allows
3rd party appliances to authenticate with it.

# Running Simple Auth

TODO

## Stand-Alone Mode

In Stand-Alone mode, you run *simple-auth* completely in isolation.  Users can sign-up on the site, and then you can authenticate them with various startegies (see below).  In this code, simple-auth doesn't communicate to anyone externally, and counts on you making an API request in the future.

### External Authentication

**Simple**: Make a `POST` to `/api/v1/auth/simple` with username/password and optionally totp parameters.  Will return 200 on access-allowed, 401 on denied.  Only works with simple-auth (eg, not external OIDC providers).  This method is least-secure, but easiest to consume in a pure-trust environment.

**Token**: TBD

**OIDC**: Full OpenID Connect OAuth2 flow.

## Same-Domain Cookie

## Reverse Proxy Gateway

# Configuration

TODO

## Login Providers

### Google

To enable google login, you need to set up an OIDC provider as documented [here](https://developers.google.com/identity/protocols/oauth2/web-server)

Then add a bit of configuration as an OIDC login-provider:

```yaml
web:
  login:
    oidc:
    - id: google
      name: 'Google'
      icon: 'fa-google'
      clientid: 'xxx'
      clientsecret: 'yyy'
      authurl: 'https://accounts.google.com/o/oauth2/v2/auth'
      tokenurl: 'https://oauth2.googleapis.com/token'
```

# Customization

Since *simple-auth* is a whitelabeled solution, it supports some level of customization via custom styles and template modifications.

The most prominant place to put styles is in `static/common.css`.  This file doesn't have anything in it by default, and can act
as a place to add overrides.

The `metadata` section of the configuration also has some pre-configured tweakable things like `company` and `copyright`.  See [simpleauth.default.yml](simpleauth.default.yml) for a full list of things to change.

# Development

## Dev-Mode

Two commands need to be run to dev:
```sh
go run simple-auth/cmd/server
npm run dev
```

## Building

```
go build -o simple-auth-server simple-auth/cmd/server
npm run build
```

OR with docker

```
docker build .
```

## Feature Wishlist

### V0
- Forgot password (single-use signin)
- Change password
- TOTP
- Account requirements (Blocks signin/usage until resolved)
  - Verification email
  - Temp ban
- UX Tweaking, autofocus, tab, enter
- Documentation
  - Docker / docker-compose gateway example
- General tests for OIDC and different cases

### V1
- OIDC Login Flow
- Google Auth
  - Better error pages
  - Refactor session to be passable and generic
  - Support JWT signature checking
- CLI tool: deactivate, activate, singleuse
  - Bundle in dockerfile
- Embed all resources into single exe?

# License

Copyright (c) 2020 Christopher LaPointe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

