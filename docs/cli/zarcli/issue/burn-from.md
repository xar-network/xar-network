# xarcli issue burn-from

## Description

某个代币的 Owner 在没有关闭持币者自己可以销毁自己持有该代币前提下，持币者对自己持有的该代币进行销毁。

## Usage
```shell
 xarcli issue burn-from [issue-id] [acc-address][amount] [flags]
```
## Flags
**Global flags, query command flags** [xarcli](../README.md)

## Example
### burn token
```shell
xarcli issue burn-from coin174876e801 xar1lgs73mwr56u2f4z4yz36w8mf7ym50e7myrqn65 88 --from $you_wallet_name
```
输入正确的密码之后，你就销毁了其他人账户里的代币。
```txt
{
  Height: 2991
  TxHash: 09E2591037100326AC7730E3E8C53103D72C1C38BF4DF82600338DD6DF38CC4B
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 29892
  Tags:
    - action = issue_burn_from
    - category = issue
    - issue-id = coin174876e802
    - sender = xar1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
}
```
