# zar unsafe-reset-all

## Description

Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state

## Usage

```shell
zar unsafe-reset-all [flags]
```

## Flags

| Name, shorthand|Default     | description               | Required  |
| ---------- | ----------- | ------------------------- | -------- |
| -h, --help |             | help for unsafe-reset-all| false  |
| --home     | ~/.zar | directory for config and data  | false    |

## Example

``` shell
zar unsafe-reset-all
```
