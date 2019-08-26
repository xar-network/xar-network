# zarcli issue describe

## Description
Owner Describes the issue token，Must be json file no larger than 1024 bytes.
## Usage
```shell
 zarcli issue describe [issue-id] [description-file] [flags]
```
## Flags

**Global flags, query command flags** [zarcli](../README.md)

## Example
### Set a description for the token
```shell
zarcli issue describe coin174876e802 /description.json --from $you_wallet_name
```
#### Template
```shell
{
    "organization":"Hashgard",
    "website":"https://www.zar.com",
    "logo":"https://cdn.zar.com/static/logo.2d949f3d.png",
    "intro":"This is a good project"
}
```
The result is as follows：
```txt
{
 Height: 3069
  TxHash: 02ED02AF5CD9C140C05D6C120BD7D57D196C27C9B3C794E6133DE912FD8243C1
  Data: 0F0E636F696E31373438373665383032
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 27465
  Tags:
    - action = issue_description
    - category = issue
    - issue-id = coin174876e802
    - sender = gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
}
```
### Query issue information
```shell
zarcli issue query-issue coin174876e802
```
The result is as follows：
```shell
{
Issue:
  IssueId:          			coin174876e802
  Issuer:           			gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
  Owner:           				gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7
  Name:             			issuename
  Symbol:    	    			AAA
  TotalSupply:      			9999999991024
  Decimals:         			18
  IssueTime:					1558179518
  Description:	    			{"org":"Hashgard","website":"https://www.zar.com","logo":"https://cdn.zar.com/static/logo.2d949f3d.png","intro":"This is a description of the project"}
  BurnOwnerDisabled:  			false
  BurnHolderDisabled:  			false
  BurnFromDisabled:  			false
  FreezeDisabled:  				false
  MintingFinished:  			false
}
```
