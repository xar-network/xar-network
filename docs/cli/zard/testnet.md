# xar testnet

## Description

Note, strict routability for addresses is turned off in the config file.

## Usage

```shell
xar testnet [flags]
```

## Flags

| Name，shorthand          | Type  | Default      | Description                                            | Required  |
| --------------------- | ------ | ------------ | ---------------------------------------------------- | -------- |
| -h, --help            |        |              | help for testnet                                    | false  |
| --chain-id            | string |              | genesis file chain-id, if left blank will be randomly created| `true`     |
| --minimum-gas-prices  | string | 0.000006gard |  Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum                      | `true`      |
| --node-cli-home       | string | xarcli  | Home directory of the node's cli configuration            | false  |
| --node-daemon-home    | string | xar     | Home directory of the node's daemon configuration| false  |
| --node-dir-prefix     | string | falsede         | Prefix the directory name for each node with (node results in node0, node1, ...) | false  |
| --output-dir          | string | ./mytestnet  | Directory to store initialization data for the testnet| false  |
| --starting-ip-address | string | 192.168.0.1  | Starting IP address                                     | false  |
| --v                   | int    | 4            |  Number of validators to initialize the testnet with| false  |

## Example

```shell
xar testnet--chain-id=${chain-id}
```
