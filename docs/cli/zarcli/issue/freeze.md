# xarcli issue freeze

## Description
When the freeze function is turned on，Owenr freezes the transfer of the specified address.
## Usage
```shell
 xarcli issue freeze [freeze-Type] [issue-id][acc-address][end-time] [flags]
```
### freeze-type

| Name   | Description            |
| ------ | -------------------- |
| in     | Transfer in|
| out    | Transfer out|
| In-out | Transfer in and out |



## Flags

**Global flags, query command flags** [xarcli](../README.md)

## Example

### Freeze specified account transfer
```shell
xarcli issue freeze in coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n\ 253382641454 --from $you_wallet_name
```
The result is as follows：
```txt
{
Height: 2570
  TxHash: DA8EEDB42B3177E281B462A88AB77D04E398286A4215D5BA0898ABA98F0270AA
  Data: 0F0E636F696E31373438373665383030
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 16459
  Tags:
    - action = issue_freeze
    - category = issue
    - issue-id = coin174876e800
    - sender = gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
    - freeze-type = in

}
```
