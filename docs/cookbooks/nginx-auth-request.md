# NGINX Authentication Request

NGINX's [auth_request](https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-subrequest-authentication/) tells nginx to make a sub-request to an external server in order to validate the user is authorized to access a resource.

<mermaid>
graph LR
A{User} -- Web Request --> B[NGINX]
B -- auth_request --> C[Simple Auth]
B -- proxy_pass --> D[Backend]
</mermaid>

## Setting up Simple-Auth with auth_request

In order to set this up, you need to do a few things:

1. Enable vouch endpoint
1. Set up nginx to sit infront of simple-auth to make an `auth_request` (vouch) to *simple-auth*
1. Set up nginx to proxy *simple-auth* UI
1. Run simple-auth server that can be proxied to by nginx

### Enabling vouch endpoint

```yaml
authenticators:
    vouch:
        enabled: true
```

### Docker

The following example will set up a docker-compose stack that has a static page sit behind simple-auth's security.

#### docker-compose

<<< @/examples/docker-nginx/docker-compose.yml

#### nginx.conf

<<< @/examples/docker-nginx/nginx.conf
