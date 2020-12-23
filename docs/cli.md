# Command Line Interface (CLI)

*simple-auth* also bundles `simple-auth-cli` which can be run as an easy interface to *simple-auth*'s database.  The CLI is able to interface with the database **without the server running**, allowing for modifications, repair, or maintenance.

::: tip
Make sure the CLI has access to your configuration via `SA_INCLUDE` environment variable, otherwise it might be pointing to the wrong (default) database.
:::

Simple actions can be completed, such as:
- Create account
- Change account password
- Create a one-time-token (Password reset)
- Examine configuration
- etc.

## Help Docs

```
NAME:
   cli - CLI Tool for simple-auth

USAGE:
   cli [global options] command [command options] [arguments...]

VERSION:
   devel, head

DESCRIPTION:
   CLI Tool for inspecting, testing, and modifying data for simple-auth

COMMANDS:
   onetime      Create one-time use token for an account
   stipulation  Modify stipulations on an account
   config       See default config
   help, h      Shows a list of commands or help for one command
   user:
     adduser  Add a new user to simple-auth DB
     passwd   Change or set password for simple-auth user

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   simple-auth  Copyright (C) 2020 Chris LaPointe
    This program comes with ABSOLUTELY NO WARRANTY.
    This is free software, and you are welcome to redistribute it
    under certain conditions
```
