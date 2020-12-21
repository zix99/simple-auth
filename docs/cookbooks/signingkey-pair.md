# SigningKey RSA Pair

Sometimes you may want to validate the JWT that has been created for session or as an OIDC token. In this case,
you need to share the secret that generated the key.

While this is easy, it's insecure.  Anyone who has the HMAC key is able to generate a new key of their own. This
is where RSA key-pairs come in.  I won't go into the details here, but you can read more about
[Public-key cryptography on Wikipedia](https://en.wikipedia.org/wiki/Public-key_cryptography).

## Setting up a key-pair

### Generating a Key

The easiest way to generate a RSA key is via `openssl`

```bash
$ openssl genrsa > privatekey.pem
```

If needed, you can extract the public key via:
```bash
$ openssl rsa -in privatekey.pem -pubout > publickey.pem
```

### Configuration

::: danger
Don't forget to insert your own key into the config below!
:::

```yaml
web:
  login:
    cookie:
      jwt:
        signingmethod: RS512
        # WARNING: Fill in your own key. Copying this from docs is insecure
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
