# xar init

## Description

Initialize validators's and node's configuration files.

## Usage

```shell
xar init [flags]
```

## Flags

| Name，shorthand| Type  | Default     | description                                                  | Required  |
| ----------- | ------ | ----------- | ------------------------------- | -------- |
| -h, --help  |        |             | help for init                             | false  |
| --chain-id  | string |             | genesis file chain-id, if left blank will be randomly created    | false  |
| --moniker   | string |             | set the validator's moniker | true    |
| --overwrite | bool   |             | overwrite the genesis.json file         | false   |
| --home      | string | ~/.xar | directory for config and data                                          | false   |
| --trace     | bool   |             |  print out full stack trace on errors                                   | false  |

## Example

`xar init --chain-id=testnet-1000 --moniker=xar`
