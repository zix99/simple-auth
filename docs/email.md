# Email

Email is optional for *simple-auth*, though it does enable some useful features like:
- Welcome email
- Email Verification
- Forgot Password

## Engines

::: warning
The default email engine is `noop`, which won't send any emails. In order to enable features like *Forgot Password*,
you must configure an engine that will send emails to the user.
:::

Simple-auth has the ability to support multiple engines, but the only one that sends out email right now is `SMTP`

```yaml
email:
  engine: noop
```

### Noop

No emails, no logging. Default.

### Stdout

Outputs the email body to the log. Useful mostly for debugging.

### SMTP

::: warning IMPORTANT
In order for the URLs in the email to function correctly, you must specify a correct `web.baseurl`.  For example, `SA_WEB_BASEURL=http://example.com`
:::

Sends email via SMTP server.

```yaml
email:
    engine: smtp
    from: null    # Who the email is "from"
    smtp:         # If engine is "smtp", the config
        host: ''         # SMTP Host
        identity: null   # Identify (often null)
        username: null   # Username
        password: null   # Password
```

For testing, you may be able to use something like [Google's SMTP Relay](https://support.google.com/a/answer/176600?hl=en)
