# zar collect-gentxs


## Description
Collect genesis txs and output a genesis.json file

## Usage
```shell
zar collect-gentxs [flags]
```

## Flags
| Name，shorthand| Type  | Default                   | description                   | Required |
| ----------- | ------ | ------------------------- | ------------------------------ | -------- |
| --gentx-dir | string | ~/.zar/config/gentx/ |  override default "gentx" directory from which collect and execute genesis transactions| false  |
| -h, --help  |        |                           |  help for collect-gentxs                    | false  |
| --home      | string | ~/.zar               |  directory for config and data              | false  |
| --trace     | bool   |                           | print out full stack trace on erro          | false  |

## Example
`zar collect-gentxs`
