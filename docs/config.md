# Config

Config can be loaded into `simple-auth` in the following ways (in-order):

1. YAML configuration file, starting with `simpleauth.default.yml` (embedded)
1. Environment configuration prefixed with `SA_`, eg `SA_PRODUCTION=true` or `SA_METADATA_COMPANY=SuperCorp`
1. Command line argument flags. eg `--metadata-company=SuperCorp`
1. Any additional configuration specified in the `include` field (see below)

## Include Config

To specify another location to look for config, you can use the `include` config (an array of paths).  This can be added in the following standard ways:

1. `include: []` in yaml config
1. `--include=/etc/simpleauth.yml`
1. `export SA_INCLUDE=/etc/simpleauth.yml`

## Default Configuration

The default config is contained purely in `simpleauth.default.yml`, which is embedded into the application executable.

<<< ../simpleauth.default.yml

# Advanced Config

## Signing Key Pair

It is possible to use a RSA key pair rather than a shared-secret style HMAC key.  Below is an example:

```yaml
web:
  login:
    cookie:
      jwt:
        signingmethod: RS512
        signingkey: |
          -----BEGIN RSA PRIVATE KEY-----
          MIICXQIBAAKBgQDei6z3v9qDh3sGe+oXgWaWa5eF6bgLpWZ+POXGpKt8DItCXrh/
          Y7qEST7UsIgEup1GLXaAL4wGbJaE3WtAdBJ/+MVDwM8FHFC/fpZ8lR6+QmXioSwg
          az2oLpq8XN/Fm59zUCRhRQ0ieZNuavto81rW4vij8w6Eu6jfz1HvI0NvRQIDAQAB
          AoGBALEmzKx2+3HiQgt6TnEhn9Ezmm2OC+SxaHIq9dn3sU5RCfXuQr2dXJb7W1mh
          oNTq3FFF1WPa9YMTo4nmW/71ptbdtrJnit/wxDauiODI6WbFonTee3h3qF08L24N
          L8InPcVBjMjO+EZL6gz43yhnROYn3LrnW1d+XqJUGwPURkiNAkEA9j8YfT6r6lYB
          F42Pn/HfQudkhtxm8K991cj+KVAnYxuczc+CY3apaJi11S8rB2KVNGBOXf9lXYp+
          1rH4FxtnFwJBAOdcQbK5Qv4DILpUpBuGCQC8B8mMWPHv0Vu+37gPL9/I37abYmI2
          iYiiEwYfEWW0TjOGWxFoz5RDZ5FRfFLc1gMCQQC7fCq/ITpvfu/x6NxToSqlm9wU
          Ojc+Rb9/SDsLZXW3pcxrfvT9mdk+RBcdq34Nb2e+qxy/wLaC0/HisTn9DeYdAkBr
          tfDEMyn0NiKTfPpV8gXY+LErbRKvuDWg3/EpsLBaIBu+2QJptPg8yy/OJsKjtdi+
          diuJpGEXpnXeGrClpzhfAkAMznhIYRsUrgdYz8b0RdYRNhXBeYj3tGYb2QRnelPb
          2BNlGlgKn8JnB/xK5gP6amGRoy0i8SKEwMcAE0QQr2Bc
          -----END RSA PRIVATE KEY-----
```