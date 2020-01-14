/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
)

// Keeper csdt Keeper
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	paramsSubspace params.Subspace
	oracle         types.OracleKeeper
	bank           types.BankKeeper
	sk             types.SupplyKeeper
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	subspace params.Subspace,
	oracle types.OracleKeeper,
	bank types.BankKeeper,
	supply types.SupplyKeeper,
) Keeper {
	return Keeper{
		storeKey:       storeKey,
		oracle:         oracle,
		bank:           bank,
		paramsSubspace: subspace.WithKeyTable(types.ParamKeyTable()),
		cdc:            cdc,
		sk:             supply,
	}
}

func (k Keeper) BuySynthetic(ctx sdk.Context, buyer sdk.AccAddress, coin sdk.Coin) sdk.Error {

	if !coin.IsValid() || coin.IsZero() || coin.IsNegative() {
		return sdk.ErrInvalidCoins("invalid coins to purchase")
	}

	// Check synthetic type ok
	p := k.GetParams(ctx)
	if !p.IsSyntheticPresent(coin.Denom) {
		return sdk.ErrInternal("synthetic type not enabled to create synthetics")
	}

	mb, found := k.GetMarketBalance(ctx, coin.Denom)
	if !found {
		snap := types.NewVolumeSnapshots(p.MarketBalanceParam.SnapshotLimit, p.MarketBalanceParam.Coefficients)
		mb = types.NewMarketBalance(coin.Denom, snap, p.MarketBalanceParam.BlocksPerSnapshot, p.MarketBalanceParam.TimerInterval)
	}

	price := k.oracle.GetCurrentPrice(ctx, coin.Denom).Price

	if price.IsNil() || price.IsZero() || price.IsNegative() {
		return sdk.ErrInternal("synthetic type does not have an oracle price")
	}

	amount := sdk.NewUintFromBigInt(price.Mul(sdk.NewDecFromInt(coin.Amount)).TruncateInt().BigInt())
	if amount.IsZero() {
		return sdk.ErrInternal("quantity too small to represent")
	}

	quantity, ok := sdk.NewIntFromString(amount.String())
	if !ok {
		return sdk.ErrInternal("quantity can not be represented")
	}
	foundersFee := mb.GetFeeForDirection(quantity, matcheng.Bid)
	quantityWithFee := p.Fee.AddToAmount(quantity).Add(foundersFee)

	purchaseCoins := sdk.NewCoins(sdk.NewCoin(types.StableDenom, quantityWithFee))
	if !purchaseCoins.IsValid() || purchaseCoins.IsAnyNegative() {
		return sdk.ErrInvalidCoins("invalid purchase coins")
	}

	ok = k.bank.HasCoins(ctx, buyer, purchaseCoins)
	if !ok {
		return sdk.ErrInsufficientCoins("not enough funds in buyer's account")
	}

	err := k.sk.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, purchaseCoins)
	if err != nil {
		panic(err) // this shouldn't happen because coin balance was checked earlier
	}

	syntheticCoins := sdk.NewCoins(coin)

	if !syntheticCoins.IsValid() || syntheticCoins.IsAnyNegative() {
		return sdk.ErrInvalidCoins("invalid synthetic coins")
	}

	er := k.sk.MintCoins(ctx, types.ModuleName, syntheticCoins)
	if er != nil {
		return er
	}

	er = k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyer, syntheticCoins)
	if er != nil {
		return er
	}

	mb.IncreaseLongVolume(quantityWithFee)
	k.SetMarketBalance(ctx, mb)
	return nil
}

func (k Keeper) SellSynthetic(ctx sdk.Context, seller sdk.AccAddress, coin sdk.Coin) sdk.Error {

	if !coin.IsValid() || coin.IsZero() || coin.IsNegative() {
		return sdk.ErrInvalidCoins("invalid coins to purchase")
	}

	// Check synthetic type ok
	p := k.GetParams(ctx)
	if !p.IsSyntheticPresent(coin.Denom) {
		return sdk.ErrInternal("synthetic type not enabled to create synthetics")
	}

	mb, found := k.GetMarketBalance(ctx, coin.Denom)
	if !found {
		snap := types.NewVolumeSnapshots(p.MarketBalanceParam.SnapshotLimit, p.MarketBalanceParam.Coefficients)
		mb = types.NewMarketBalance(coin.Denom, snap, p.MarketBalanceParam.BlocksPerSnapshot, p.MarketBalanceParam.TimerInterval)
	}

	price := k.oracle.GetCurrentPrice(ctx, coin.Denom).Price

	if price.IsNil() || price.IsZero() || price.IsNegative() {
		return sdk.ErrInternal("synthetic type does not have an oracle price")
	}

	amount := sdk.NewUintFromBigInt(price.Mul(sdk.NewDecFromInt(coin.Amount)).TruncateInt().BigInt())
	if amount.IsZero() {
		return sdk.ErrInternal("quantity too small to represent")
	}

	quantity, ok := sdk.NewIntFromString(amount.String())
	if !ok {
		return sdk.ErrInternal("quantity can not be represented")
	}

	foundersFee := mb.GetFeeForDirection(quantity, matcheng.Ask)
	quantitySubFee := p.Fee.SubFromAmount(quantity).Sub(foundersFee)

	syntheticCoins := sdk.NewCoins(coin)
	if !syntheticCoins.IsValid() || syntheticCoins.IsAnyNegative() {
		return sdk.ErrInvalidCoins("invalid synthetic coins")
	}

	ok = k.bank.HasCoins(ctx, seller, syntheticCoins)
	if !ok {
		return sdk.ErrInsufficientCoins("not enough funds in seller's account")
	}

	err := k.sk.SendCoinsFromAccountToModule(ctx, seller, types.ModuleName, syntheticCoins)
	if err != nil {
		panic(err) // this shouldn't happen because coin balance was checked earlier
	}

	sellerCoins := sdk.NewCoins(sdk.NewCoin(types.StableDenom, quantitySubFee))

	if !sellerCoins.IsValid() || sellerCoins.IsAnyNegative() {
		return sdk.ErrInvalidCoins("invalid seller coins")
	}

	er := k.sk.BurnCoins(ctx, types.ModuleName, syntheticCoins)
	if er != nil {
		return er
	}

	er = k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, seller, sellerCoins)
	if er != nil {
		return er
	}

	mb.IncreaseShortVolume(quantitySubFee)
	k.SetMarketBalance(ctx, mb)
	return nil
}

func (k Keeper) GetStableDenom() string {
	return types.StableDenom
}
func (k Keeper) GetGovDenom() string {
	return types.GovDenom
}

// GetOracle allows testing
func (k Keeper) GetOracle() types.OracleKeeper {
	return k.oracle
}

// GetOracle allows testing
func (k Keeper) GetSupply() types.SupplyKeeper {
	return k.sk
}

func (k Keeper) IsNominee(ctx sdk.Context, nominee string) bool {
	params := k.GetParams(ctx)
	nominees := params.Nominees
	for _, v := range nominees {
		if v == nominee {
			return true
		}
	}
	return false
}
