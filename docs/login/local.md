# Local (Username)

Simple authentication provides a mechanism for simple-auth to store username and password in its own database and allow login via its various UI and API mechanisms.

## Configuration

### Creation

Various configuration for account-creation exist, including the enablement of it.

```yaml
providers:
    settings:
        createaccountenabled: true      # If allowed to create account
    local:
        emailvalidationrequired: false  # If email validation is required before login
```

::: warning
In order for email validation to work, email must be enabled. See [email](../cookbooks/email)
:::

### Requirements

Requirements allow enforcing username/password characters, strength, and length.

```yaml
providers:
    local:
        requirements:
            usernameregex: '^[a-z][a-z0-9.]+$'
            passwordminlength: 3
            passwordmaxlength: 30
            usernameminlength: 5
            usernamemaxlength: 20
```

## Features

### reCAPTCHA v2

[reCAPTCHA](https://www.google.com/recaptcha/about/) is a way for *simple-auth* to validate that the user creating or logging into an account is not a bot.  We use reCAPTCHA v2.

Technically, this is a feature of the web-interface, but often comes in handy when using
local authentication.

If enabled, the recaptcha prompt will show on the create user page, and forgot password.

::: tip
To setup, you first need a site key from [google's recaptcha service](https://developers.google.com/recaptcha/intro)
:::

```yaml
web:
    recaptchav2:
        enabled: false
        sitekey: null    # site key from google
        secret: null     # Secret key from google
        theme: 'light'   # light or dark
```

### TOTP (2FA)

Two-factor authentication (2FA) prompts the user for a code from a device, in addition to a password, to allow them to login.

TOTP presents the user with both the QR Code and the secret. Most popular apps should function fine (eg. Authy, or Google Authenticator)

```yaml
providers:
    local:
        twofactor:
            enabled: true
            issuer: "simple-auth"
```

### Forgot Password

::: warning
By default, forgot password isn't enabled because it relies on a [email engine](../cookbooks/email) being set up.
:::

*Forgot password* functionality can send an an email to an account with a link to
update their password.  When an account is logged-into via this one-time token,
they will be able to change their password without knowing the original.

Once you have email setup, you can enable forgot-password by enabling the following config.

```yaml
web:
    login:
        onetime:
            allowforgotpassword: true
```
