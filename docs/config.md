# Config

Config can be loaded into `simple-auth` in the following ways (in-order):

1. YAML configuration file, starting with `simpleauth.default.yml` (embedded)
1. Environment configuration prefixed with `SA_`, eg `SA_PRODUCTION=true` or `SA_METADATA_COMPANY=SuperCorp`
1. Command line argument flags. eg `--metadata-company=SuperCorp`
1. Any additional configuration specified in the `include` field (see below)

::: tip
When translating from *yaml* to other formats (eg. `metadata.company`):
1. When an argument, seprate with `-`, eg `--metadata-company`
1. When an environment variable, separate with a `_`, eg `SA_METADATA_COMPANY`
:::

::: tip
For boolean flags, when setting via argument, you can shorthand it by not providing a value (eg `--web-prometheus`)
:::

## Include Config

To specify another location to look for a **yaml config file**, you can use the `include` config.  This can be added in the following standard ways:

1. `include: []` in yaml config
1. `--include=/etc/simpleauth.yml`
1. `export SA_INCLUDE=/etc/simpleauth.yml`

## Default Configuration

The default config is contained purely in `simpleauth.default.yml`, which is embedded into the application executable.  This file contains all of the defaults for the application.

<<< ../simpleauth.default.yml

# Advanced Config

## Persistence

If you want to use a database other than `sqlite3` (local file), see [database drivers](/database)

## Signing Key Pair

See [SigningKey RSA Pair Cookbook](cookbooks/signingkey-pair)

