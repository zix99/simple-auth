# simple-auth

Simple-auth is a lightweight, whitelabeled, authentication solution.  It allows users to sign-up
with a UI, offering various security measures.  Then, via various authenticators, allows
3rd party appliances to authenticate with it.

## Development

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

