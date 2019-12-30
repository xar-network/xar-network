/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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
	"bytes"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// Keeper csdt Keeper
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             *codec.Codec
	paramsSubspace  params.Subspace
	oracle          types.OracleKeeper
	bank            types.BankKeeper
	sk              types.SupplyKeeper
	liquidityModule string
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	subspace params.Subspace,
	oracle types.OracleKeeper,
	bank types.BankKeeper,
	supply types.SupplyKeeper,
	liquidityModule string,
) Keeper {
	return Keeper{
		storeKey:        storeKey,
		oracle:          oracle,
		bank:            bank,
		paramsSubspace:  subspace.WithKeyTable(types.ParamKeyTable()),
		cdc:             cdc,
		sk:              supply,
		liquidityModule: liquidityModule,
	}
}

// ModifyCSDT creates, changes, or deletes a CSDT
// TODO can/should this function be split up?
func (k Keeper) ModifyCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInCollateral sdk.Int, changeInDebt sdk.Int) sdk.Error {

	// Check the owner has enough collateral and stable coins
	err := validateCoinTransfer(ctx, k, owner, collateralDenom, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.changeState(ctx, owner, collateralDenom, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.executeCoinTransfer(ctx, owner, collateralDenom, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) executeCoinTransfer(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInCollateral, changeInDebt sdk.Int) sdk.Error {
	err := k.handleCollateralChange(ctx, owner, collateralDenom, changeInCollateral)
	if err != nil {
		return err
	}

	err = k.handleDebtChange(ctx, owner, collateralDenom, changeInDebt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) changeState(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInCollateral, changeInDebt sdk.Int) sdk.Error {

	csdt := k.createOrGetCSDT(ctx, owner, collateralDenom)

	err := k.changeCsdtState(ctx, &csdt, collateralDenom, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.validateAndSetGlobalDebt(ctx, changeInDebt)
	if err != nil {
		return err
	}

	err = k.validateAndSetCollateralState(ctx, changeInDebt, &csdt)
	if err != nil {
		return err
	}

	if csdt.CollateralAmount.IsZero() && csdt.Debt.IsZero() { // TODO maybe abstract this logic into SetCSDT
		k.DeleteCSDT(ctx, csdt)
	} else {
		k.SetCSDT(ctx, csdt)
	}
	return err
}

func (k Keeper) handleDebtChange(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInDebt sdk.Int) sdk.Error {
	p := k.GetParams(ctx)
	if changeInDebt.IsNegative() { //Depositing stable coin from owner to CSDT (decrease supply)
		depositCoins := sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt.Neg()))

		err := k.sk.SendCoinsFromAccountToModule(ctx, owner, k.liquidityModule, depositCoins)
		if err != nil {
			return err
		}

		return k.sk.BurnCoins(ctx, k.liquidityModule, depositCoins)
	}
	//Withdrawing stable coins to owner (minting)
	stableCoin := p.Fee.AddToCoin(sdk.NewCoin(types.StableDenom, changeInDebt))
	withdrawCoins := sdk.NewCoins(stableCoin)

	err := k.sk.MintCoins(ctx, k.liquidityModule, withdrawCoins)
	if err != nil {
		return err
	}

	return k.sk.SendCoinsFromModuleToAccount(ctx, k.liquidityModule, owner, withdrawCoins)
}

func (k Keeper) handleCollateralChange(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInCollateral sdk.Int) sdk.Error {
	// change owner's coins (increase or decrease)
	if changeInCollateral.IsNegative() {
		return k.sk.SendCoinsFromModuleToAccount(ctx, k.liquidityModule, owner, sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral.Neg())))
	}

	return k.sk.SendCoinsFromAccountToModule(ctx, owner, k.liquidityModule, sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral)))
}

// TODO: should we handle a case when csdt == nil?
func (k Keeper) changeCsdtState(ctx sdk.Context, csdt *types.CSDT, collateralDenom string, changeInCollateral, changeInDebt sdk.Int, ) sdk.Error {
	err := addCollateralToCsdt(collateralDenom, changeInCollateral, csdt)
	if err != nil {
		return err
	}

	err = addDebtToCsdt(changeInDebt, csdt)
	if err != nil {
		return err
	}

	currentPrice := k.oracle.GetCurrentPrice(ctx, csdt.CollateralDenom).Price
	liquidationRatio := k.oracle.GetCurrentPrice(ctx, csdt.CollateralDenom).Price
	err = csdt.Validate(currentPrice, liquidationRatio)
	if err != nil {
		return err
	}

	return err
}

func addCollateralToCsdt(collateralDenom string, changeInCollateral sdk.Int, csdt *types.CSDT) sdk.Error {
	var collateralCoins sdk.Coins

	if changeInCollateral.IsNegative() {
		collateralCoins = sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral.Neg()))
		csdt.CollateralAmount = csdt.CollateralAmount.Sub(collateralCoins)

	} else {
		collateralCoins = sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral))
		csdt.CollateralAmount = csdt.CollateralAmount.Add(collateralCoins)
	}
	return nil
}

func addDebtToCsdt(changeInDebt sdk.Int, csdt *types.CSDT) sdk.Error {
	var debtCoins sdk.Coins

	if csdt.CollateralAmount.IsAnyNegative() {
		return sdk.ErrInternal(" can't withdraw more collateral than exists in CSDT")
	}

	if changeInDebt.IsNegative() {
		debtCoins = sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt.Neg()))
		csdt.Debt = csdt.Debt.Sub(debtCoins)
	} else {
		debtCoins = sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt))
		csdt.Debt = csdt.Debt.Add(debtCoins)
	}

	if csdt.Debt.IsAnyNegative() {
		return sdk.ErrInternal("can't pay back more debt than exists in CSDT")
	}

	return nil
}

func validateCoinTransfer(ctx sdk.Context, k Keeper, owner sdk.AccAddress, collateralDenom string, changeInCollateral sdk.Int, changeInDebt sdk.Int) sdk.Error {
	p := k.GetParams(ctx)
	if !p.IsCollateralPresent(collateralDenom) { // maybe abstract this logic into GetCSDT
		return sdk.ErrInternal("collateral type not enabled to create CSDTs")
	}

	if changeInCollateral.IsPositive() { // adding collateral to CSDT
		ok := k.bank.HasCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral)))
		if !ok {
			return sdk.ErrInsufficientCoins("not enough collateral in sender's account")
		}
	}
	if changeInDebt.IsNegative() { // reducing debt, by adding stable coin to CSDT
		ok := k.bank.HasCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt.Neg())))
		if !ok {
			return sdk.ErrInsufficientCoins("not enough stable coin in sender's account")
		}
	}
	return nil
}

// TODO
// // TransferCSDT allows people to transfer ownership of their CSDTs to others
// func (k Keeper) TransferCSDT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, collateralDenom string) sdk.Error {
// 	return nil
// }

// PartialSeizeCSDT removes collateral and debt from a CSDT and decrements global debt counters. It does not move collateral to another account so is unsafe.
// TODO should this be made safer by moving collateral to liquidatorModuleAccount ? If so how should debt be moved?
func (k Keeper) PartialSeizeCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, collateralToSeize sdk.Int, debtToSeize sdk.Int) sdk.Error {
	// get CSDT
	csdt, found := k.GetCSDT(ctx, owner, collateralDenom)
	if !found {
		return sdk.ErrInternal("could not find CSDT")
	}

	// Check if CSDT is undercollateralized
	p := k.GetParams(ctx)
	isUnderCollateralized := csdt.IsUnderCollateralized(
		k.oracle.GetCurrentPrice(ctx, csdt.CollateralDenom).Price,
		p.GetCollateralParam(csdt.CollateralDenom).LiquidationRatio,
	)
	if !isUnderCollateralized {
		return sdk.ErrInternal("CSDT is not currently under the liquidation ratio")
	}

	// Remove Collateral
	if collateralToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative collateral")
	}
	collateralCoins := sdk.NewCoins(sdk.NewCoin(csdt.CollateralDenom, collateralToSeize))
	csdt.CollateralAmount = csdt.CollateralAmount.Sub(collateralCoins)
	if csdt.CollateralAmount.IsAnyNegative() {
		return sdk.ErrInternal("can't seize more collateral than exists in CSDT")
	}

	// Remove Debt
	if debtToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative debt")
	}
	debtCoins := sdk.NewCoins(sdk.NewCoin(types.StableDenom, debtToSeize))
	csdt.Debt = csdt.Debt.Sub(debtCoins)
	if csdt.Debt.IsAnyNegative() {
		return sdk.ErrInternal("can't seize more debt than exists in CSDT")
	}

	// Update debt per collateral type
	collateralState, found := k.GetCollateralState(ctx, csdt.CollateralDenom)
	if !found {
		return sdk.ErrInternal("could not find collateral state")
	}
	collateralState.TotalDebt = collateralState.TotalDebt.Sub(debtToSeize)
	if collateralState.TotalDebt.IsNegative() {
		return sdk.ErrInternal("Total debt per collateral type is negative.") // This should not happen given the checks on the CSDT.
	}

	// Note: Global debt is not decremented here. It's only decremented when debt and stable coin are annihilated (aka heal)
	// TODO update global seized debt? this is what maker does (named vice in Vat.grab) but it's not used anywhere

	// Store updated state
	if csdt.CollateralAmount.IsZero() && csdt.Debt.IsZero() { // TODO maybe abstract this logic into SetCSDT
		k.DeleteCSDT(ctx, csdt)
	} else {
		k.SetCSDT(ctx, csdt)
	}
	k.SetCollateralState(ctx, collateralState)
	return nil
}

// ReduceGlobalDebt decreases the stored global debt counter. It is used by the liquidator when it annihilates debt and stable coin.
// TODO Can the interface between csdt and liquidator modules be improved so that this function doesn't exist?
func (k Keeper) ReduceGlobalDebt(ctx sdk.Context, amount sdk.Int) sdk.Error {
	if amount.IsNegative() {
		return sdk.ErrInternal("reduction in global debt must be a positive amount")
	}
	newGDebt := k.GetGlobalDebt(ctx).Sub(amount)
	if newGDebt.IsNegative() {
		return sdk.ErrInternal("cannot reduce global debt by amount specified")
	}
	k.SetGlobalDebt(ctx, newGDebt)
	return nil
}

func (k Keeper) GetStableDenom() string {
	return types.StableDenom
}
func (k Keeper) GetGovDenom() string {
	return types.GovDenom
}

// ---------- Store Wrappers ----------

func (k Keeper) getCSDTKeyPrefix(collateralDenom string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte("csdt"),
			[]byte(collateralDenom),
		},
		nil, // no separator
	)
}
func (k Keeper) getCSDTKey(owner sdk.AccAddress, collateralDenom string) []byte {
	return bytes.Join(
		[][]byte{
			k.getCSDTKeyPrefix(collateralDenom),
			[]byte(owner.String()),
		},
		nil, // no separator
	)
}
func (k Keeper) GetCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string) (types.CSDT, bool) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// get CSDT
	bz := store.Get(k.getCSDTKey(owner, collateralDenom))
	// unmarshal
	if bz == nil {
		return types.CSDT{}, false
	}
	var csdt types.CSDT
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &csdt)
	return csdt, true
}

func (k Keeper) createOrGetCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string) types.CSDT {
	csdt, found := k.GetCSDT(ctx, owner, collateralDenom)
	if !found {
		csdt = types.CSDT{
			Owner:            owner,
			CollateralDenom:  collateralDenom,
			CollateralAmount: sdk.NewCoins(sdk.NewCoin(collateralDenom, sdk.ZeroInt())),
			Debt:             sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.ZeroInt())),
			AccumulatedFees:  sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.ZeroInt())),
		}
		return csdt
	}
	return csdt
}

//Potentially change this logic to use the account interface?
func (k Keeper) SetCSDT(ctx sdk.Context, csdt types.CSDT) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(csdt)
	store.Set(k.getCSDTKey(csdt.Owner, csdt.CollateralDenom), bz)
}
func (k Keeper) DeleteCSDT(ctx sdk.Context, csdt types.CSDT) { // TODO should this id the csdt by passing in owner,collateralDenom pair?
	// get store
	store := ctx.KVStore(k.storeKey)
	// delete key
	store.Delete(k.getCSDTKey(csdt.Owner, csdt.CollateralDenom))
}

// GetCSDTs returns all CSDTs, optionally filtered by collateral type and liquidation price.
// `price` filters for CSDTs that will be below the liquidation ratio when the collateral is at that specified price.
func (k Keeper) GetCSDTs(ctx sdk.Context, collateralDenom string, price sdk.Dec) (types.CSDTs, sdk.Error) {
	// Validate inputs
	p := k.GetParams(ctx)
	if len(collateralDenom) != 0 && !p.IsCollateralPresent(collateralDenom) {
		return nil, sdk.ErrInternal("collateral denom not authorized")
	}
	if len(collateralDenom) == 0 && !(price.IsNil() || price.IsNegative()) {
		return nil, sdk.ErrInternal("cannot specify price without collateral denom")
	}

	// Get an iterator over CSDTs
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, k.getCSDTKeyPrefix(collateralDenom)) // could be all CSDTs is collateralDenom is ""

	// Decode CSDTs into slice
	var csdts types.CSDTs
	for ; iter.Valid(); iter.Next() {
		var csdt types.CSDT
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &csdt)
		csdts = append(csdts, csdt)
	}

	// Sort by collateral ratio (collateral/debt)
	sort.Sort(types.ByCollateralRatio(csdts)) // TODO this doesn't make much sense across different collateral types

	// Filter for CSDTs that would be under-collateralized at the specified price
	// If price is nil or -ve, skip the filtering as it would return all CSDTs anyway
	if !price.IsNil() && !price.IsNegative() {
		var filteredCSDTs types.CSDTs
		for _, csdt := range csdts {
			if csdt.IsUnderCollateralized(price, p.GetCollateralParam(collateralDenom).LiquidationRatio) {
				filteredCSDTs = append(filteredCSDTs, csdt)
			} else {
				break // break early because list is sorted
			}
		}
		csdts = filteredCSDTs
	}

	return csdts, nil
}

var globalDebtKey = []byte("globalDebt")

func (k Keeper) GetGlobalDebt(ctx sdk.Context) sdk.Int {
	// get store
	store := ctx.KVStore(k.storeKey)
	// get bytes
	bz := store.Get(globalDebtKey)
	// unmarshal
	if bz == nil {
		panic("global debt not found")
	}
	var globalDebt sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &globalDebt)
	return globalDebt
}

// Add/Subtract from global debt limit
func (k Keeper) validateAndSetGlobalDebt(ctx sdk.Context, changeInDebt sdk.Int) sdk.Error {
	p := k.GetParams(ctx)
	gDebt := k.GetGlobalDebt(ctx)
	gDebt = gDebt.Add(changeInDebt)
	if gDebt.IsNegative() {
		return sdk.ErrInternal("global debt can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if gDebt.GT(p.GlobalDebtLimit.AmountOf(types.StableDenom)) {
		return sdk.ErrInternal("change to CSDT would put the system over the global debt limit")
	}
	k.SetGlobalDebt(ctx, gDebt)
	return nil
}

// Add/Subtract from collateral debt limit
func (k Keeper) validateAndSetCollateralState(ctx sdk.Context, changeInDebt sdk.Int, csdt *types.CSDT) sdk.Error {
	p := k.GetParams(ctx)
	collateralState, found := k.GetCollateralState(ctx, csdt.CollateralDenom)
	if !found {
		collateralState = types.CollateralState{Denom: csdt.CollateralDenom, TotalDebt: sdk.ZeroInt()} // Already checked that this denom is authorized, so ok to create new CollateralState
	}
	collateralState.TotalDebt = collateralState.TotalDebt.Add(changeInDebt)
	if collateralState.TotalDebt.IsNegative() {
		return sdk.ErrInternal("total debt for this collateral type can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if collateralState.TotalDebt.GT(p.GetCollateralParam(csdt.CollateralDenom).DebtLimit.AmountOf(types.StableDenom)) {
		return sdk.ErrInternal("change to CSDT would put the system over the debt limit for this collateral type")
	}
	k.SetCollateralState(ctx, collateralState)
	return nil
}

func (k Keeper) SetGlobalDebt(ctx sdk.Context, globalDebt sdk.Int) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(globalDebt)
	store.Set(globalDebtKey, bz)
}

func (k Keeper) getCollateralStateKey(collateralDenom string) []byte {
	return []byte(collateralDenom)
}
func (k Keeper) GetCollateralState(ctx sdk.Context, collateralDenom string) (types.CollateralState, bool) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// get bytes
	bz := store.Get(k.getCollateralStateKey(collateralDenom))
	// unmarshal
	if bz == nil {
		return types.CollateralState{}, false
	}
	var collateralState types.CollateralState
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &collateralState)
	return collateralState, true
}
func (k Keeper) SetCollateralState(ctx sdk.Context, collateralstate types.CollateralState) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(collateralstate)
	store.Set(k.getCollateralStateKey(collateralstate.Denom), bz)
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
