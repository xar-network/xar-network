# xar unsafe-reset-all

## Description

Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state

## Usage

```shell
xar unsafe-reset-all [flags]
```

## Flags

| Name, shorthand|Default     | description               | Required  |
| ---------- | ----------- | ------------------------- | -------- |
| -h, --help |             | help for unsafe-reset-all| false  |
| --home     | ~/.xar | directory for config and data  | false    |

## Example

``` shell
xar unsafe-reset-all
```
