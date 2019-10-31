# xar start

## Description

Run the full node

## Usage

```shell
xar start [flags]
```

## Flags

| Name，shorthand                    | Default               | description              | Required  |
| ------------------------------- | --------------------- | ------------------------------ | -- |
| --abci                          | socket                | Specify abci transport (socket or grpc) | false      |
| --address                       | tcp://0.0.0.0:26658   | Listen address               | false  |
| --consensus.create_empty_blocks | true                  |  Set this to false to only produce blocks when there are txs or when the AppHash changes  | false  |
| --fast_sync                     | true                  | Fast blockchain syncing                            | false  |
| --minimum_fees                  |                       |  Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum                            | false  |
| --moniker                       | instance-c5m0fg87     | falsede Name         | false  |
| --p2p.laddr                     | tcp://0.0.0.0:26656   |  Node listen address. (0.0.0.0:0 means any interface, any port)   | false  |
| --p2p.persistent_peers          |                       | Comma-delimited ID@host:port persistent peers   | false  |
| --p2p.pex                       | true                  | Enable/disable Peer-Exchange       | false  |
| --p2p.private_peer_ids          |                       | Comma-delimited private peer IDs       | false  |
| --p2p.seed_mode                 |                       | Enable/disable seed mode          | false  |
| --p2p.seeds                     |                       | Comma-delimited ID@host:port seed nodes     | false  |
| --p2p.upnp                      |                       | Enable/disable UPNP port forwarding                     | false  |
| --priv_validator_laddr          |                       | Socket address to listen on for connections from external priv_validator process| false  |
| --proxy_app                     | tcp://127.0.0.1:26658 | Proxy app address, or one of: 'kvstore', 'persistent_kvstore', 'counter', 'counter_serial' or 'noop' for local testing. | false  |
| --pruning                       | syncable              | Pruning strategy: syncable, nothing, everything     | false  |
| --replay                        |                       | Replay the last block                                  | false  |
| --rpc.grpc_laddr                |                       | GRPC listen address (BroadcastTx only). Port required | false  |
| --rpc.laddr                     | tcp://0.0.0.0:26657   |  RPC listen address. Port required            | false  |
| --rpc.unsafe                    |                       | Enabled unsafe rpc methods                              | false  |
| --trace-store                   |                       | Enable KVStore tracing to an output file                | false  |
| --with-tendermint               | true                  | Run abci app embedded in-process with tendermint | false  |
| -h, --help                      |                       | help for start                           | false  |

## Example

```shell
xar start --home=/root/.xar
```
