# zarlcd User documentation

## Basic function introduction

1. Provide restful APIs and swagger-ui to show these APIs
2. Verify query proof


## zarlcd Usage

zarlcd has three subcommands:

| subcommand  | Description           | Example command                           |
| ------- | ------------------------- | ----------------------------------------- |
| help    | Print the zarlcd help           | zarlcd help               |
| version | Print the zarlcd version        | zarlcd version            |
| start   | Start zarLCD node | zarlcd start --chain-id=`<chain-id>` |

### start subcommands

`start` subcommand has these options:

#### Flags

| Parameter   | Type  | Default                 | Required| description                         |
| ------------ | ------ | ----------------------- | -------- | ----- |
| cors         | string |   | false    | Set the domains that can make CORS requests（\*all）   |
| laddr     | string | "tcp://localhost:1317"  | false  | Address for server to listen on|
| max-open     | int    | 1000                    | false    | The number of maximum open connections   |
| ssl-certfile | int    | 1000                    | false    | SSl certificate directory，Not set will automatically generate a new certificate  |
| ssl-hosts    | int    | 1000                    | false    | Domain name using ssl certificate            |
| ssl-keyfile  | int    | 1000                    | false    | ssl-key 文件所在目录，不设置 ssl 证书时会被忽略 |
| tls          | bool   | false                   | false    | Open SSL/TLS                                    |
| height       | int    | （last block height）    | false    | Get the latest block                          |
| falsede         | string | "tcp://localhost:26657" | false    | Full node rpc address                            |
| trust-node   | bool   | false                   | false    | Trust connected full nodes (Don't verify proofs for responses)   |
| help         | bool   | false                   | false    | Print the help                                |
| indent       | bool   | false                   | false    | Output result formatting                      |
| ledger       | bool   | false                   | false    | Use ledger  wallet                            |

#### Flags

| Parameter | Type  | Default               | Required| description                                      |
| -------- | ------ | --------------------- | -------- | ----------------------------------------------- |
| chain-id | string | null                  | true     | Chain ID of tendermint node                      |
| encoding | string | "hex"                 | false    | Binary encoding (hex|b64|btc)|
| home     | string | "\$HOME/.zarlcd" | false    | directory for config and data|
| output   | string | text                  | false    | Output format (text|json)  |
| trace    | bool   | false                 | false    | print out full stack trace on errors                     |

## Example Commands

1. By default, zarLCD doesn't trust the connected full node. But if you are sure about that the connected full node is trustable, then you should run zarLCD with --trust-node option:

```shell
zarlcd start --chain-id=<chain-id> --trust-node
```

2. If you want to access your zarlcd in another machine, you have to specify --laddr, for instance:

```shell
zarlcd start --chain-id=<chain-id> --laddr=tcp://0.0.0.0:1317
```

## REST APIs

Once zarlcd is started,  you can open localhost:1317/swagger-ui/ in your explorer and all restful APIs will be shown. The swagger-ui page has detailed description about APIs' functionality and required parameters. Here we just list all APIs and introduce their functionality briefly.

1. Tendermint APIs, such as query blocks, transactions and validator set
    1. `GET /node_info`: The properties of the connected node
    2. `GET /syncing`: Syncing state of node
    3. `GET /blocks/latest`: Get the latest block
    4. `GET /blocks/{height}`: Get a block at a certain height
    5. `GET /validatorsets/latest`:  Get the latest validator set
    6. `GET /validatorsets/{height}`: Get a validator set at a certain height
    7. `GET /txs/{hash}`:  Get a Tx by hash
    8. `GET /txs`: Search transactions
    9. `POST /txs`: Broadcast transaction
2. Key management APIs
    1. `GET /keys`: List of accounts stored locally
    2. `POST /keys`:  Create a new account locally
    3. `GET /keys/seed`: Create a new seed to create a new account with
    4. `GET /keys/{name}`:  Get a certain locally stored account
    5. `PUT /keys/{name}`: Update the password for this account in the KMS
    6. `DELETE /keys/{name}`: Remove an account
    7. `GET /auth/accounts/{address}`: Query information about the key object account
3. Sign and broadcast transactions APIs
    1. `POST /tx/sign`:  Sign a transaction
    2. `POST /tx/broadcast`: Broadcast a signed StdTx with amino encoding signature and public key
    3. `POST /txs/send`: Send a signed StdTx with amino encoding signature and public key
    4. `GET /bank/coin/{coin-Type}`: Get coin Type
    5. `GET /bank/token-stats`: Get token statistic
    6. `GET /bank/balances/{address}`: Get account token balances
    7. `POST /bank/accounts/{address}/transfers`: Send coins
    8. `POST /bank/burn`: burn token
4. Stake module APIs
    1. `POST /stake/delegators/{delegatorAddr}/delegate`:  Submit delegation transaction
    2. `POST /stake/delegators/{delegatorAddr}/redelegate`: Submit redelegation transaction
    3. `POST /stake/delegators/{delegatorAddr}/unbond`: Submit unbonding transaction
    4. `GET /stake/delegators/{delegatorAddr}/delegations`: Get all delegations from a delegator
    5. `GET /stake/delegators/{delegatorAddr}/unbonding_delegations`: Get all unbonding delegations from a delegator
    6. `GET /stake/delegators/{delegatorAddr}/redelegations`: Get all redelegations from a delegator
    7. `GET /stake/delegators/{delegatorAddr}/validators`: Query all validators that a delegator is bonded to
    8. `GET /stake/delegators/{delegatorAddr}/validators/{validatorAddr}`: Query a validator that a delegator is bonded to
    9. `GET /stake/delegators/{delegatorAddr}/txs`: Get all staking txs from a delegator
    10. `GET /stake/delegators/{delegatorAddr}/delegations/{validatorAddr}`: Query the current delegation between a delegator and a validator
    11. `GET /stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}`: Query all unbonding delegations between a delegator and a validator
    12. `GET /stake/validators`: Get all validator candidates
    13. `GET /stake/validators/{validatorAddr}`: Query the information from a single validator
    14. `GET /stake/validators/{validatorAddr}/unbonding_delegations`:Get all unbonding delegations from a validator
    15. `GET /stake/validators/{validatorAddr}/redelegations`: Get all outgoing redelegations from a validator
    16. `GET /stake/pool`: Get the current state of the staking pool
    17. `GET /stake/parameters`: Get the current staking parameter values
5. Governance module APIs
    1. `POST /gov/proposal`: Submit a proposal
    2. `GET /gov/proposals`: Query proposals
    3. `POST /gov/proposals/{proposalId}/deposits`: Deposit tokens to a proposal
    4. `GET /gov/proposals/{proposalId}/deposits`: Query deposits
    5. `POST /gov/proposals/{proposalId}/votes`: Vote a proposal
    6. `GET /gov/proposals/{proposalId}/votes`: Query voters
    7. `GET /gov/proposals/{proposalId}`:  Query a proposal
    8. `GET /gov/proposals/{proposalId}/deposits/{depositor}`:Query deposit
    9. `GET /gov/proposals/{proposalId}/votes/{voter}`: Query vote
    10. `GET/gov/params`: Query governance parameters
6. Slashing module APIs
    1. `GET /slashing/validators/{validatorPubKey}/signing_info`: Get sign info of given validator
    2. `POST /slashing/validators/{validatorAddr}/unjail`:  Unjail a jailed validator
7. Distribution module APIs
    1. `POST /distribution/{delegatorAddr}/withdrawAddress`: Set withdraw address
    2. `GET /distribution/{delegatorAddr}/withdrawAddress`: Get withdraw address
    3. `POST /distribution/{delegatorAddr}/withdrawReward`: Withdraw reward
    4. `GET /distribution/{delegatorAddr}/distrInfo/{validatorAddr}`: Get the revenue distribution information of a delegate
    5. `GET /distribution/{delegatorAddr}/distrInfos`: Query the entrusted income distribution information of all the principals
    6. `GET /distribution/{validatorAddr}/valDistrInfo`: Query the validator revenue distribution information
8. Query app version
    1. `GET /version`: Version of zarHUB
    2. `GET /node_version`: Version of the connected node

## Special Parameters

These apis are picked out from above section. And they can be used to build and broadcast transactions:

1. `POST /bank/accounts/{address}/transfers`
2. `POST /stake/delegators/{delegatorAddr}/delegate`
3. `POST /stake/delegators/{delegatorAddr}/redelegate`
4. `POST /stake/delegators/{delegatorAddr}/unbond`
5. `POST /gov/proposal`
6. `POST /gov/proposals/{proposalId}/deposits`
7. `POST /gov/proposals/{proposalId}/votes`
8. `POST /slashing/validators/{validatorAddr}/unjail`

They all support the these special query parameters below. By default, their values are all false. And each parameter has its unique priority( Here 0 is the top priority). If multiple parameters are specified to true, then the parameters with lower priority will be ignored. For instance, if generate-only is true, then all other parameters, such as simulate and commit will be ignored.

| parameter name| Type| Default| Priority| Description                 |
| ------------- | ---- | ------ | ------ | -------------------------- |
| generate-only | bool | false  | 0      | Build an unsigned transaction and return it back|
| simulate      | bool | false  | 1      | Ignore the gas field and perform a simulation of a transaction, but don’t broadcast it  |
| commit        | bool | false  | 2      | Wait for transaction being included in a block |
| async         | bool | false  | 3      | Broadcast transaction asynchronously |
