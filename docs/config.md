# Config

Config can be loaded into `simple-auth` in the following ways (in-order):

1. YAML configuration file, starting with `simpleauth.default.yml`, which will then load `simpleauth.yml` by default, and recurse with any other `include` statements
1. Environment configuration prefixed with `SA_`, eg `SA_PRODUCTION=true` or `SA_METADATA_COMPANY=SuperCorp`
1. Command line argument flags. eg `--metadata-company=SuperCorp`

## Nested Config

By default, *simple-auth* will look for a `simpleauth.yml` file located in the working directory.  If it's not found, it will present a **WARN**, but otherwise continue as normal.

If you want to specify another location to look for config, you can use the `include` config (an array of paths).  This can be added in the following standard ways:

1. `include: []` in yaml config
1. `--include=/etc/simpleauth.yml`
1. `export SA_INCLUDE=/etc/simpleauth.yml`

## Default Configuration

The default config is contained purely in `simpleauth.default.yml`, which is boxed into the application by default.

<<< @/simpleauth.default.yml
