# Quick-Start

[[toc]]

## Stand-Alone Authentication Portal

The easiest way to get started is to run *simple-auth* as a stand-alone portal, accessed by API, to manage users.

By default, only one config value is required to be set: `web.login.cookie.jwt.signingkey`  This key must be unique, and cryptographically secure, to guarantee a user's session can't be hijacked.

### Docker

```sh
docker run -it --rm -e SA_WEB_LOGIN_COOKIE_JWT_SIGNINGKEY=REPLACE_ME zix99/simple-auth
```

### Binary

Download the binary from the github releases page, and run with:

```sh
./simple-auth-server --web-login-cookie-jwt-signingkey=REPLACE_ME
```

## Next Steps
