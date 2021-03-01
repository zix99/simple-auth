# LetsEncrypt TLS

*simple-auth* provides the ability to automatically issue a valid TLS certificate by leveraging [Let's Encrypt](https://letsencrypt.org/).

## Enabling

To enable, you simply need to set `web.tls.enabled` to `true`.

For added security, you can provide a list of hostnames that we're allowed to issue a certificate for via `web.tls.autohosts`.

## Config

```yaml
web:
    tls:
        enabled: true
        certfile: null   # Certificate file, if enabled (and not auto)
        keyfile: null    # Key file, if enabled (and not auto)
        # AutoTLS (and cache) are used to leverage LetsEncrypt to acquire certificate
        # Needs to be internet-facing to work
        auto: true       # If false, will use certfile and keyfile instead of letsencrypt
        autohosts: []    # Optional list of hosts that we're allowed to issue a cert for
        cache: ./
```
