# zarcli issue search

## Description
Search issues based on symbol

## Usage
```shell
zarcli issue search [symbol] [flags]
```
## Flags

**Global flags, query command flags** [zarcli](../README.md)

## Example

### Search
```shell
zarcli issue search AAA
```
```txt
 [
    {
        "issue_id":"coin174876e802",
        "issuer":"gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7",
        "owner":"gard1f76ncl7d9aeq2thj98pyveee8twplfqy3q4yv7",
        "issue_time":"1558179518",
        "name":"issuename",
        "symbol":"AAA",
        "total_supply":"10000000001023",
        "decimals":"18",
        "description":"{"org":"Hashgard","website":"https://www.zar.com","logo":"https://cdn.zar.com/static/logo.2d949f3d.png","intro":"This is a description of the project"}",
        "burn_owner_disabled":false,
        "burn_holder_disabled":false,
        "burn_from_disabled":false,
        "freeze_disabled":false,
        "minting_finished":false
    }
]

```
