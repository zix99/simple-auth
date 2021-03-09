# TLS (HTTPS)

There are primarily two ways to enable https (TLS) on *simple-auth*: Certificates, and Let's Encrypt

[[toc]]

## Let's Encrypt TLS

*simple-auth* provides the ability to automatically issue a valid TLS certificate by leveraging [Let's Encrypt](https://letsencrypt.org/).

### Enabling

::: tip Prerequisites
*simple-auth* needs to be exposed to the public internet, and have a domain, in order to obtain a certificate from Let's Encrypt
:::

To enable, you simply need to set `web.tls.enabled` to `true`.

For added security, you can provide a list of hostnames that we're allowed to issue a certificate for via `web.tls.autohosts`.

#### How does it work?

When a user first accesses *simple-auth*, if there is no certificate, then it will automatically make a call to LetsEncrypt with
the correct callback url.  If the host is on the `autohosts` list (or that list is empty), a certificate will be issued, cached, and
then used to secure the connection going forward.

### Config

::: tip
In docker, the default cache directory will be in the same volume as the DB
:::

```yaml
web:
    tls:
        enabled: true
        # AutoTLS (and cache) are used to leverage LetsEncrypt to acquire certificate
        # Needs to be internet-facing to work
        auto: true       # If false, will use certfile and keyfile instead of letsencrypt
        autohosts: []    # Optional list of hosts that we're allowed to issue a cert for
        cache: ./tlscache
```

## Certificates

### Getting SSL Certificate

The following command will create a **self-signed** certificate you can use for *simple-auth*.  This certificate **will not be recognized as valid by the browser** unless you create and install your own certificate authority.  Alternatively, you may obtain a valid certificate from a certificate authority.

```bash
openssl req -x509 -nodes -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
```

### Config

```yaml
web:
    tls:
        enabled: true
        certfile: null   # Certificate file, if enabled (and not auto)
        keyfile: null    # Key file, if enabled (and not auto)
        auto: false      # Need to disable Let's Encrypt
```
