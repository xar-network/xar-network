# xarcli issue mint

## Description
With the additional switch turned on，Owner Add tokens for yourself or add tokens to others。

## Usage
```shell
 xarcli issue mint [issue-id] [amount] [flags]
```
| Name   | Type    | Required   | Default   | Description      |
| --------  | ------------------- | ----- | ------ | -------- |
| --to                  | string | false|| Add tokens to the specified address |

## Flags

**Global flags, query command flags** [xarcli](../README.md)

## Example

### Add tokens to the specified address
```shell
xarcli issue mint coin174876e802 9999 --to=xar1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7 --from $you_wallet_name
```
The result is as follows：
```txt
{
  Height: 3138
  TxHash: 110F99B71B2F206E29EDA2A5EC9DB1E372045693C06EDB9C32B9C9767AB92F93
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 40402
  Tags:
    - action = issue_mint
    - category = issue
    - issue-id = coin174876e802
    - sender = xar1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
}
```
