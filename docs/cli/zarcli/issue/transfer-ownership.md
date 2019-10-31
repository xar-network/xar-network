# xarcli issue transfer-ownership

## Description
Token owner transfer the ownership to new account
## Usage
```shell
 xarcli issue transfer-ownership [issue-id] [to_address] [flags]
```
## Flags

**Global flags, query command flags** [xarcli](../README.md)

## Example
### transfer ownership
```shell
 xarcli issue transfer-ownership coin174876e802 gard1lgs73mwr56u2f4z4yz36w8mf7ym50e7myrqn65 --from $you_wallet_name
```
The result is as follows：
```txt
{
   Height: 3199
  TxHash: 3438C2C4F054730CD02FC30C408B3DA558CE9C5CC99810F83406DB1D41708CC9
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 26680
  Tags:
    - action = issue_transfer_ownership
    - category = issue
    - issue-id = coin174876e802
    - sender = gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
}
```
