# xar tendermint

## Description

Tendermint subcommands

## Usage

```shell
xar tendermint [subcommand] [flags]
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
| --home     | string | ~/.xar | directory for config and data         | false  |
| --trace    | bool   |             | print out full stack trace on errors | false  |

## Example

```shell
xar tendermint show-node-id
xar tendermint show-validator
xar tendermint show-address
```
