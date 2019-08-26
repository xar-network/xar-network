# zar export

## Description

zar can export blockchain state at any height and output json format string.

## Usage

```shell
zar export [flags]
```

## Flags

| Name，shorthand      | Type  | Default| description                                 | Required  |
| ----------------- | ------ | ------ | ------------------------------------------- | -------- |
| -h, --help        |        |        | help for export                          | false   |
| --for-zero-height |        |        | Export state to start at height zero   | false   |
| --height          | int    | -1     | Export state from a particular height   | false  |
| --jail-whitelist  | string |        | List of validators to not jail state export| false  |

## Example

`zar export`
