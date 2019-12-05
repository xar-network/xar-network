# CSDT Specification

## Overview

3 problem statements we wanted to solve;

1. Staking is inflationary
2. Rewards are volatile (can't price server cost accordingly)
3. Can't issue AUM on-chain of $200MM if the staking token staked is worth $2MM (51% attacks become too profitable)

So we created a staking mechanism (called CSDT) that means, new FTM isn't minted, rewards are in USD, and you can use multi collateral to stake

You use FTM (or any other collateral as voted for by FTM token holders) to create ucsdt

You stake ucsdt and receive rewards in ucsdt

ucsdt can be settled to unlock FTM, or traded on the DEX or atomic swap modules

This way, supply remains fixed, rewards are stable, people can (optionally) swap ucsdt back to uftm

## CSDT Version 1 (hub_1 implementation)

```go
// CSDT is the state of a single account.
type CSDT struct {
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`                        
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	CollateralAmount sdk.Int        `json:"collateral_amount" yaml:"collateral_amount"`
	Debt             sdk.Int        `json:"debt" yaml:"debt"`
}
```

```
owner: bech32 address
collateral_denom: string
collateral_amount: int64
debt: int64
```

CSDT version 1 was designed to store a CSDT on a per owner, per denomination relationship. So for each denomination (BTC/ETH/FTM) you would have a CSDT. Version 1 CSDTs could only mint ucsdt.

Collateral is subtracted from an owners balance, essentially erasing the collateral from existence.

## CSDT Version 2 (hub_2 implementation)

```go
// CSDT is the state of a single account.
type CSDT struct {
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	CollateralAmount sdk.Coins      `json:"collateral_amount" yaml:"collateral_amount"`
	Debt             sdk.Coins      `json:"debt" yaml:"debt"`
}
```

```
owner: bech32 address
collateral: [{denom:string,amount:int64}]
debt:[{denom:ucsdt:string,amount:int64}]
```

Version 2 saves CSDTs per owner and allow for an array of collateral to be used. This collateral creates a overall collateral value in ucsdt. The formula being;

```
foreach collateral in collateral_array
- value += amount*(priceOf(denom))
```

This creates an overall basket ucsdt value and this value of ucsdt can be minted.

With Version 2, collateral is transferred from the from the owner into the CSDT module.

## CSDT Version 3 (hub_3 upgrade)

```go
// CSDT is the state of a single account.
type CSDT struct {
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`
	CollateralAmount sdk.Coins      `json:"collateral_amount" yaml:"collateral_amount"`
	Debt             sdk.Coins      `json:"debt" yaml:"debt"`
}
```

```
owner: bech32 address
collateral: [{denom:string,amount:int64}]
debt:[{denom:string,amount:int64}]
```

Version 2 saves CSDTs per owner and allow for an array of collateral to be used. This collateral creates a overall collateral value in ucsdt. The formula being;

```
foreach collateral in collateral_array
- value += amount*(priceOf(denom))
```

With the version 3 upgrade, the CSDT module now has an amount of collateral. Should ucsdt be requested, this is first borrowed from the CSDT module, before it is minted.

With version 3, debt also becomes a basket, and an owner can borrow any asset that the CSDT module owns (access to all other collateral). So at this upgrade, an owner could borrow BTC/ETH against FTM/ucsdt (including adding ucsdt as collateral)

## CSDT Version 4 (hub_4 upgrade)

```go
// CSDT is the state of a single account.
type CSDT struct {
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`
	CollateralAmount sdk.Coins      `json:"collateral_amount" yaml:"collateral_amount"`
	Debt             sdk.Coins      `json:"debt" yaml:"debt"`
	AccumulatedFees  sdk.Coins      `json:"accumulated_fees" yaml:"accumulated_fees"`
}
```

```
owner: bech32 address
collateral: [{denom:string,amount:int64}]
debt:[{denom:string,amount:int64}]
fees: [{denom:string,amount:int64}]
```

At hub_4 upgrade borrowing and minting will accrue fees. The fee will be based on the amount of liquidity provided vs the borrowed capacity. 


U = Borrowings / (Cash + Borrowings)
Borrowing Interest Rate = 2% + (U x 20%)
Lending Interest Rate = Borrowing Interest Rate x U

Fees are calculated per block. If ucsdt is borrowed, then the interest is based off of the calculation above, if ucsdt is minted, then it incurs a stability fee that will be set via on-chain governance. ~1% at genesis.

## CSDT Version 4.5 (hub_4 upgrade)

With the hub_4 upgrade the CSDT module is overcollateralized thanks to interest. This interest will be paid out to all overcollateralized CSDTs. The parameter at genesis will be set to 200% but can be controlled via on-chain governance.

Meaning all accounts that have more than 200% collateral, will receive equal portions of fees for the collateral they provide (augmented for the interest bearing calculation above)
