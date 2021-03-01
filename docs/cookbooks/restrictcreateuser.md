# Restrict Create User

It is possible to disable the web create-user form in order to restrict who can login via *simple-auth*.  Once disabled, the only way to create-user manually will be via manual intervention (the CLI).

## Disabling Web Create User

To disable create-user, set `providers.settings.createaccountenabled` to `false`. This will disable the links, web-flow, and APIs to create a user in *simple-auth*.

## Creating User on CLI

### Binary

To create a user, run `adduser` from the CLI.  It will read any configuration the same way *simple-auth* does, and create the user in the database.

```sh
./simple-auth-cli adduser "Firstname lastname" "email@example.com" "username"
```

::: tip
The password will be prompted for, or can be provided as the last argument after the username.
:::

### Docker

```sh
docker exec -it <container-id> ./simple-auth-cli adduser "Firstname lastname" "email@example.com" "username"
```

::: tip
The password will be prompted for, or can be provided as the last argument after the username.
:::
