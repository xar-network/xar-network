# zarcli issue create

## Description

Issue a new token

## Usage

```shell
zarcli issue create [name] [symbol] [total-supply] [flags]
```

## Flags

| Name          | Type| Required  | Default| Description                               |
| ------------- | ---- | -------- | ------ | --------------------------------------- |
| --decimals    | int  | true   | 18     | Decimals of the token |
| --burn-owner  | Bool | false   | false  | Disable owner to burn your own token|
| --burn-holder | bool | false  | false  | Disable Non-owner users burn their own tokens|
| --burn-from   | bool | false   | false  | Disable owner burning other user tokens|
| --minting     | bool | false   | false  | Disable the mint              |
| --freeze      | bool | false | false  | Disable freeze          |

**Global flags, query command flags** [zarcli](../README.md)

## Example

### Issue a new coin

```shell
zarcli issue create issuename AAA 10000000000000 --from $you_wallet_name
```

The result is as follows：

```txt
{
   Height: 2967
  TxHash: 84B19F831958A6334C4806967E66E6C8640F0A2E7958A5E99A1DF3B6B6E6378C
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 43428
  Tags:
    - action = issue
    - category = issue
    - issue-id = coin174876e802
    - sender = gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
}
```

Query account

```shell
zarcli bank account gard1f203m5q7hr4tkf0vredrn4wpxkx7zngn4pntye
```

There is a `coin (issue-id)` token in your token list.

```shell
{
 Account:
  Address:       gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
  Pubkey:        gardpub1addwnpepqfpd8mkl3jg43fw7y02fe99cgaxutf5npv9y9gx9dvrrcdwl36shv694apw
  Coins:         9999999990001issuename(coin174876e802)
  AccountNumber: 0
  Sequence:      16
}
```
