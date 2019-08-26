# zar tendermint

## Description

Tendermint subcommands

## Usage

```shell
zar tendermint [subcommand] [flags]
```

## Subcommands

| Name          | Description                            |
| ---------------- | ------------------------------------ |
| --show-node-id   | Show this node's ID                   |
| --show-validator | Show this node's tendermint validator info |
| --show-address   | Shows this node's tendermint validator consensus address |

## Flags

| Name, shorthand|Type  | Default     | description              | Required  |
| ---------- | ------ | ----------- | ------------------------ | -------- |
| -h, --help |        |             | help for tendermint       | false  |
| --home     | string | ~/.zar | directory for config and data         | false  |
| --trace    | bool   |             | print out full stack trace on errors | false  |

## Example

```shell
zar tendermint show-node-id
zar tendermint show-validator
zar tendermint show-address
```
