# xarcli issue burn

## Description
Token holder or the Owner burn the token he holds
## Usage
```shell
 xarcli issue burn [issue-id] [amount] [flags]
```
## Flags

**Global flags, query command flags** [xarcli](../README.md)

## Example
### burn coin
```shell
xarcli issue burn coin174876e800 88888 --from $you_wallet_name
```
The result is as follows：
```txt
{
   Height: 3020
  TxHash: 9C74FB0071940687E026EDEAB3666F8E3C0624C8541ABCF61C6BBFBFBA533F97
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 27544
  Tags:
    - action = issue_burn_holder
    - category = issue
    - issue-id = coin174876e802
    - sender = gard1lgs73mwr56u2f4z4yz36w8mf7ym50e7myrqn65
}
```
