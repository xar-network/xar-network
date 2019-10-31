# xard add-genesis-account

## Description
add genesis account to /path/to/.xard/config/genesis.json


## Usage
```shell
xard add-genesis-account [address_or_key_name] [coin][,[coin]] [flags]
```


## Subcommands
| Name         | Type  | Default| description                | Required |
| --------------------- | ------ | ------ | ------------------- | -------- |
| [address_or_key_name] | string |        | Added account name or address    | true    |
| [coin]                | string |        | coin type and amount | true    |


## Flags
| Name，shorthand         | Type  | Default        | Description                      | Required |
| -------------------- | ------ | -------------- | -------------------------------- | -------- |
| -h, --help           |        |                | help for add-genesis-account  | false  |
| --home-client        | string | ~/.xarcli | client's home directory       | false   |
| --vesting-amount     | string |                | amount of coins for vesting accounts  | false    |
| --vesting-end-time   | int    |                | schedule end time (unix epoch) for vesting accounts| false    |
| --vesting-start-time | int    |                | schedule start time (unix epoch) for vesting accounts| false    |
| --home               | string | ~/.xard    | directory for config and data| false    |
| --trace              | bool   |                | print out full stack trace on errors| false   |


## Example
```shell
xarcli keys add root
xard add-genesis-account root 100000000xar
```
