package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetMarketInfo(ctx sdk.Context, name string) MoneyMarket {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewMarket()
	}
	bz := store.Get([]byte(name))
	var moneymarket MoneyMarket
	k.cdc.MustUnmarshalBinaryBare(bz, &moneymarket)
	return moneymarket
}

// HasMarketOwner - returns whether or not the Market already has an owner
func (k Keeper) HasMarketOwner(ctx sdk.Context, name string) bool {
	return !k.GetMarketInfo(ctx, name).Owner.Empty()
}

// GetMarketOwner - get the current owner of a Market
func (k Keeper) GetMarketOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetMarketInfo(ctx, name).Owner
}

// Sets the entire Market metadata struct for a name
func (k Keeper) SetMarketInfo(ctx sdk.Context, name string, moneymarket MoneyMarket) {
	if moneymarket.Owner.Empty() {
		return
	}
	fmt.Println(moneymarket)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(moneymarket))
}

// SetMarketOwner - sets the current owner of a market
func (k Keeper) SetMarketOwner(ctx sdk.Context, name string, symbol string, owner sdk.AccAddress, tokenName string, collateralToken string) {
	fmt.Println(name)
	moneymarket := k.GetMarketInfo(ctx, name)
	moneymarket.Owner = owner
	moneymarket.Name = name
	moneymarket.Symbol = symbol
	moneymarket.TokenName = tokenName
	moneymarket.CollateralToken = collateralToken
	k.SetMarketInfo(ctx, name, moneymarket)
}

//SetMarketOwner - sets the current owner of a market
func (k Keeper) SupplyMarketPosition(ctx sdk.Context, Owner sdk.AccAddress, market string, lendTokens sdk.Coins) {
	if Owner.Empty() {
		return
	}
	marketposition := k.GetMarketPosition(ctx, Owner)
	marketposition.Owner = Owner
	marketposition.Market = market
	marketposition.LendTokens = marketposition.LendTokens.Add(lendTokens)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(Owner), k.cdc.MustMarshalBinaryBare(marketposition))

	//set the market
	moneymarket := k.GetMarketInfo(ctx, market)
	moneymarket.TokenSupply = moneymarket.TokenSupply.Add(lendTokens)
	interestRate := calcInterest(moneymarket.TokenBorrows, moneymarket.TokenSupply)
	moneymarket.InterestRate = interestRate
	k.SetMarketInfo(ctx, market, moneymarket)
}

// SetMarketOwner - sets the current owner of a market
func (k Keeper) GetMarketPosition(ctx sdk.Context, Owner sdk.AccAddress) MarketPosition {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(Owner)) {
		return NewMarketPosition()
	}
	bz := store.Get([]byte(Owner))
	var marketposition MarketPosition
	k.cdc.MustUnmarshalBinaryBare(bz, &marketposition)
	return marketposition
}

// BorrowFromMarketPosition - borrow from a token market
func (k Keeper) BorrowFromMarketPosition(ctx sdk.Context, Owner sdk.AccAddress, market string, borrowTokens sdk.Coins, borrowCollateral sdk.Coins) {
	if Owner.Empty() {
		return
	}
	marketposition := k.GetMarketPosition(ctx, Owner)
	marketposition.Owner = Owner
	marketposition.Market = market
	marketposition.BorrowTokens = marketposition.BorrowTokens.Add(borrowTokens)
	marketposition.BorrowCollateral = marketposition.BorrowCollateral.Add(borrowCollateral)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(Owner), k.cdc.MustMarshalBinaryBare(marketposition))

	//set the market
	moneymarket := k.GetMarketInfo(ctx, market)
	moneymarket.TokenBorrows = moneymarket.TokenBorrows.Add(borrowTokens)
	moneymarket.BorrowCollateral = moneymarket.BorrowCollateral.Add(borrowCollateral)
	fmt.Println(moneymarket)
	interestRate := calcInterest(moneymarket.TokenBorrows, moneymarket.TokenSupply)
	moneymarket.InterestRate = interestRate
	k.SetMarketInfo(ctx, market, moneymarket)
}

// RedeemFromMarketPosition - redeem from a token market
func (k Keeper) RedeemFromMarketPosition(ctx sdk.Context, Owner sdk.AccAddress, market string, redeemTokens sdk.Coins) {
	if Owner.Empty() {
		return
	}
	marketposition := k.GetMarketPosition(ctx, Owner)
	marketposition.Owner = Owner
	marketposition.Market = market
	marketposition.LendTokens = marketposition.LendTokens.Sub(redeemTokens)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(Owner), k.cdc.MustMarshalBinaryBare(marketposition))

	//set the market
	moneymarket := k.GetMarketInfo(ctx, market)
	moneymarket.TokenSupply = moneymarket.TokenSupply.Sub(redeemTokens)
	interestRate := calcInterest(moneymarket.TokenBorrows, moneymarket.TokenSupply)
	moneymarket.InterestRate = interestRate
	k.SetMarketInfo(ctx, market, moneymarket)
}

// RedeemFromMarketPosition - redeem from a token market
func (k Keeper) RepayToMarketPosition(ctx sdk.Context, Owner sdk.AccAddress, market string, repayTokens sdk.Coins, borrowCollateral sdk.Coins) {
	if Owner.Empty() {
		return
	}
	marketposition := k.GetMarketPosition(ctx, Owner)
	marketposition.Owner = Owner
	marketposition.Market = market
	marketposition.BorrowTokens = marketposition.BorrowTokens.Sub(repayTokens)
	marketposition.BorrowCollateral = marketposition.BorrowCollateral.Sub(borrowCollateral)
	//fmt.Println(marketposition)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(Owner), k.cdc.MustMarshalBinaryBare(marketposition))

	//set the market
	// moneymarket := k.GetMarketInfo(ctx, market)
	// moneymarket.TokenBorrows = moneymarket.TokenBorrows.Sub(repayTokens)
	// fmt.Println("printing market info")
	// fmt.Println(moneymarket)
	// interestRate := calcInterest(moneymarket.TokenBorrows, moneymarket.TokenSupply)
	// moneymarket.InterestRate = interestRate
	// k.SetMarketInfo(ctx, market, moneymarket)
}

func calcInterest(tokenBorrows sdk.Coins, tokenSupply sdk.Coins) sdk.Int {
	fmt.Println(tokenBorrows[0].Amount)
	allMoney := tokenBorrows.Add(tokenSupply)
	fmt.Println(allMoney[0].Amount)
	interestRate := tokenBorrows[0].Amount.Quo(allMoney[0].Amount)
	fmt.Println(interestRate)
	//interestRate := 2
	return interestRate
}

// NewKeeper creates new instances of the moneymarket Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}
