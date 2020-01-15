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
func (k Keeper) ModifyCSDT(ctx sdk.Context, owner sdk.AccAddress, changeInCollateral sdk.Coin, changeInDebt sdk.Coin) sdk.Error {

	// Check the owner has enough collateral and stable coins
	err := validateCoinTransfer(ctx, k, owner, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.changeState(ctx, owner, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.executeCoinTransfer(ctx, owner, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) executeCoinTransfer(ctx sdk.Context, owner sdk.AccAddress, changeInCollateral, changeInDebt sdk.Coin) sdk.Error {
	err := k.handleCollateralChange(ctx, owner, changeInCollateral)
	if err != nil {
		return err
	}

	err = k.handleDebtChange(ctx, owner, changeInDebt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) changeState(ctx sdk.Context, owner sdk.AccAddress, changeInCollateral sdk.Coin, changeInDebt sdk.Coin) sdk.Error {

	csdt := k.createOrGetCSDT(ctx, owner, changeInCollateral, changeInDebt)

	err := k.changeCsdtState(ctx, &csdt, changeInCollateral, changeInDebt)
	if err != nil {
		return err
	}

	err = k.validateAndSetGlobalDebt(ctx, changeInDebt)
	if err != nil {
		return err
	}

	err = k.validateAndSetCollateralState(ctx, changeInDebt)
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

func (k Keeper) handleDebtChange(ctx sdk.Context, owner sdk.AccAddress, changeInDebt sdk.Coin) sdk.Error {
	p := k.GetParams(ctx)
	if changeInDebt.IsNegative() { //Depositing stable coin from owner to CSDT (decrease supply)
		depositCoins := sdk.NewCoins(sdk.NewCoin(changeInDebt.Denom, changeInDebt.Amount.Neg()))

		err := k.sk.SendCoinsFromAccountToModule(ctx, owner, k.liquidityModule, depositCoins)
		if err != nil {
			return err
		}

		if changeInDebt.Denom == types.StableDenom {
			// If stabelDenom in debt - use old logic with Burn
			return k.sk.BurnCoins(ctx, k.liquidityModule, depositCoins)
		}

		return nil
	}
	//Withdrawing stable coins to owner (minting)
	withdrawCoins := sdk.NewCoins(p.Fee.AddToCoin(changeInDebt))

	if changeInDebt.Denom == types.StableDenom {
		// If stabelDenom in debt - use old logic with Mint
		err := k.sk.MintCoins(ctx, k.liquidityModule, withdrawCoins)
		if err != nil {
			return err
		}
	}

	return k.sk.SendCoinsFromModuleToAccount(ctx, k.liquidityModule, owner, withdrawCoins)
}

func (k Keeper) handleCollateralChange(ctx sdk.Context, owner sdk.AccAddress, changeInCollateral sdk.Coin) sdk.Error {
	// change owner's coins (increase or decrease)
	if changeInCollateral.IsNegative() {
		return k.sk.SendCoinsFromModuleToAccount(ctx, k.liquidityModule, owner, sdk.NewCoins(sdk.NewCoin(changeInCollateral.Denom, changeInCollateral.Amount.Neg())))
	}

	return k.sk.SendCoinsFromAccountToModule(ctx, owner, k.liquidityModule, sdk.NewCoins(changeInCollateral))
}

// TODO: should we handle a case when csdt == nil?
func (k Keeper) changeCsdtState(ctx sdk.Context, csdt *types.CSDT, changeInCollateral, changeInDebt sdk.Coin, ) sdk.Error {
	err := addCollateralToCsdt(changeInCollateral, csdt)
	if err != nil {
		return err
	}

	err = addDebtToCsdt(changeInDebt, csdt)
	if err != nil {
		return err
	}

	// TODO: Is require other checks for csdt.DebtDenom?
	currentPrice := k.oracle.GetCurrentPrice(ctx, changeInCollateral.Denom).Price
	liquidationRatio := k.oracle.GetCurrentPrice(ctx, changeInCollateral.Denom).Price
	err = csdt.Validate(currentPrice, liquidationRatio, changeInCollateral.Denom)
	if err != nil {
		return err
	}

	return err
}

func addCollateralToCsdt(changeInCollateral sdk.Coin, csdt *types.CSDT) sdk.Error {
	var collateralCoins sdk.Coins

	if changeInCollateral.IsNegative() {
		collateralCoins = sdk.NewCoins(sdk.NewCoin(changeInCollateral.Denom, changeInCollateral.Amount.Neg()))
		csdt.CollateralAmount = csdt.CollateralAmount.Sub(collateralCoins)

	} else {
		collateralCoins = sdk.NewCoins(changeInCollateral)
		csdt.CollateralAmount = csdt.CollateralAmount.Add(collateralCoins)
	}
	return nil
}

func addDebtToCsdt(changeInDebt sdk.Coin, csdt *types.CSDT) sdk.Error {
	var debtCoins sdk.Coins

	if csdt.CollateralAmount.IsAnyNegative() {
		return sdk.ErrInternal(" can't withdraw more collateral than exists in CSDT")
	}

	if changeInDebt.IsNegative() {
		debtCoins = sdk.NewCoins(sdk.NewCoin(changeInDebt.Denom, changeInDebt.Amount.Neg()))
		csdt.Debt = csdt.Debt.Sub(debtCoins)
	} else {
		debtCoins = sdk.NewCoins(changeInDebt)
		csdt.Debt = csdt.Debt.Add(debtCoins)
	}

	if csdt.Debt.IsAnyNegative() {
		return sdk.ErrInternal("can't pay back more debt than exists in CSDT")
	}

	return nil
}

func validateCoinTransfer(ctx sdk.Context, k Keeper, owner sdk.AccAddress, changeInCollateral sdk.Coin, changeInDebt sdk.Coin) sdk.Error {
	p := k.GetParams(ctx)
	if !p.IsCollateralPresent(changeInCollateral.Denom) { // maybe abstract this logic into GetCSDT
		return sdk.ErrInternal("collateral type not enabled to create CSDTs")
	}

	if changeInCollateral.IsPositive() { // adding collateral to CSDT
		ok := k.bank.HasCoins(ctx, owner, sdk.NewCoins(changeInCollateral))
		if !ok {
			return sdk.ErrInsufficientCoins("not enough collateral in sender's account")
		}
	} else {
		// Check decrease limitations
		cParm := p.GetCollateralParam(changeInCollateral.Denom)

		if !k.isValidDecreaseLimits(ctx, cParm) {
			return sdk.ErrInsufficientCoins("not enough collateral in global pool account (try later)")
		}
	}
	if changeInDebt.IsNegative() { // reducing debt, by adding coin to CSDT
		ok := k.bank.HasCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(changeInDebt.Denom, changeInDebt.Amount.Neg())))
		if !ok {
			return sdk.ErrInsufficientCoins("not enough stable coin in sender's account")
		}
	}
	return nil
}

func (k Keeper) GetPoolValue(ctx sdk.Context, poolDenom string) sdk.Int {
	return k.sk.GetModuleAccount(ctx, k.liquidityModule).GetCoins().AmountOf(poolDenom)
}

func (k Keeper) isValidDecreaseLimits(ctx sdk.Context, parm types.CollateralParam) bool {
	// Check all decrease limits
	if parm.DecreaseLimits == nil {
		return true
	}

	poolVal := k.GetPoolValue(ctx, parm.Denom)
	snap, _ := k.GetPoolSnapshot(ctx)

	for _, lim := range parm.DecreaseLimits {
		// Get snapshot and current pool value
		snapCoin := snap.GetVal(lim, parm.Denom)
		if snapCoin == nil {
			// If have not snapshot for this denom and period - ignore this limit
			continue
		}
		snapVal := snapCoin.Amount

		// Min pool value for current limit
		borderLimit := snapVal.Sub(snapVal.Mul(sdk.NewInt(100)).Mod(lim.MaxPercent))
		if poolVal.LTE(borderLimit) {
			// If pool value lower then min pool value from limit - return false
			return false
		}
	}

	// If not detect any limits - return true ("isValid")
	return true
}

// TODO
// // TransferCSDT allows people to transfer ownership of their CSDTs to others
// func (k Keeper) TransferCSDT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, collateralDenom string) sdk.Error {
// 	return nil
// }

// PartialSeizeCSDT removes collateral and debt from a CSDT and decrements global debt counters. It does not move collateral to another account so is unsafe.
// TODO should this be made safer by moving collateral to liquidatorModuleAccount ? If so how should debt be moved?
func (k Keeper) PartialSeizeCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, collateralToSeize sdk.Int, debtDenom string, debtToSeize sdk.Int) sdk.Error {
	// get CSDT
	csdt, found := k.GetCSDT(ctx, owner)
	if !found {
		return sdk.ErrInternal("could not find CSDT")
	}

	// Check if CSDT is undercollateralized
	p := k.GetParams(ctx)
	isUnderCollateralized := csdt.IsUnderCollateralized(
		k.oracle.GetCurrentPrice(ctx, collateralDenom).Price,
		p.GetCollateralParam(collateralDenom).LiquidationRatio,
		collateralDenom,
	)
	if !isUnderCollateralized {
		return sdk.ErrInternal("CSDT is not currently under the liquidation ratio")
	}

	// Remove Collateral
	if collateralToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative collateral")
	}
	collateralCoins := sdk.NewCoins(sdk.NewCoin(collateralDenom, collateralToSeize))
	csdt.CollateralAmount = csdt.CollateralAmount.Sub(collateralCoins)
	if csdt.CollateralAmount.IsAnyNegative() {
		return sdk.ErrInternal("can't seize more collateral than exists in CSDT")
	}

	// Remove Debt
	if debtToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative debt")
	}
	debtCoins := sdk.NewCoins(sdk.NewCoin(debtDenom, debtToSeize))
	csdt.Debt = csdt.Debt.Sub(debtCoins)
	if csdt.Debt.IsAnyNegative() {
		return sdk.ErrInternal("can't seize more debt than exists in CSDT")
	}

	// Update debt per collateral type
	collateralState, found := k.GetCollateralState(ctx, collateralDenom)
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

func (k Keeper) GetStableDenom() string {
	return types.StableDenom
}
func (k Keeper) GetGovDenom() string {
	return types.GovDenom
}

// ---------- Store Wrappers ----------

func (k Keeper) getCSDTKeyPrefix(denom string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte("csdt"),
			[]byte(denom),
		},
		nil, // no separator
	)
}
func (k Keeper) getCSDTKey(owner sdk.AccAddress) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(owner.String()),
		},
		nil, // no separator
	)
}
func (k Keeper) GetCSDT(ctx sdk.Context, owner sdk.AccAddress) (types.CSDT, bool) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// get CSDT
	bz := store.Get(k.getCSDTKey(owner))
	// unmarshal
	if bz == nil {
		return types.CSDT{}, false
	}
	var csdt types.CSDT
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &csdt)
	return csdt, true
}

func (k Keeper) createOrGetCSDT(ctx sdk.Context, owner sdk.AccAddress, changeInCollateral, changeInDebt sdk.Coin) types.CSDT {
	csdt, found := k.GetCSDT(ctx, owner)
	if !found {
		csdt = types.CSDT{
			Owner:            owner,
			CollateralAmount: sdk.NewCoins(sdk.NewCoin(changeInCollateral.Denom, sdk.ZeroInt())),
			Debt:             sdk.NewCoins(sdk.NewCoin(changeInDebt.Denom, sdk.ZeroInt())),
			AccumulatedFees:  sdk.NewCoins(sdk.NewCoin(changeInDebt.Denom, sdk.ZeroInt())),
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
	store.Set(k.getCSDTKey(csdt.Owner), bz)
}
func (k Keeper) DeleteCSDT(ctx sdk.Context, csdt types.CSDT) { // TODO should this id the csdt by passing in owner,collateralDenom pair?
	// get store
	store := ctx.KVStore(k.storeKey)
	// delete key
	store.Delete(k.getCSDTKey(csdt.Owner))
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
			if csdt.IsUnderCollateralized(price, p.GetCollateralParam(collateralDenom).LiquidationRatio, collateralDenom) {
				filteredCSDTs = append(filteredCSDTs, csdt)
			} else {
				break // break early because list is sorted
			}
		}
		csdts = filteredCSDTs
	}

	return csdts, nil
}

// Pool snapshots
func (k Keeper) getPoolSnapshotKey() []byte {
	return bytes.Join(
		[][]byte{
			[]byte("poolsnapshot"),
		},
		nil, // no separator
	)
}
func (k Keeper) GetPoolSnapshot(ctx sdk.Context) (types.PoolSnapshot, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(k.getPoolSnapshotKey())
	if bz == nil {
		return types.PoolSnapshot{}, false
	}
	var snap types.PoolSnapshot
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &snap)
	return snap, true
}

func (k Keeper) createOrGetPoolSnapshot(ctx sdk.Context) types.PoolSnapshot {
	snap, found := k.GetPoolSnapshot(ctx)
	if !found {
		snap = types.PoolSnapshot{
			ByLimits: make([]types.PoolSnapValue, 0),
		}
		return snap
	}
	return snap
}
func (k Keeper) SetPoolSnapshot(ctx sdk.Context, snap types.PoolSnapshot) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(snap)
	store.Set(k.getPoolSnapshotKey(), bz)
}

// Add/Subtract from global debt limit
// TODO: How validate GlobalDebt with debtDenom?
func (k Keeper) validateAndSetGlobalDebt(ctx sdk.Context, changeInDebt sdk.Coin) sdk.Error {
	p := k.GetParams(ctx)
	collateralState, found := k.GetCollateralState(ctx, changeInDebt.Denom)
	if !found {
		collateralState = types.CollateralState{Denom: changeInDebt.Denom, TotalDebt: sdk.ZeroInt()} // Already checked that this denom is authorized, so ok to create new CollateralState
	}
	collateralState.TotalDebt = collateralState.TotalDebt.Add(changeInDebt.Amount)
	if collateralState.TotalDebt.IsNegative() {
		return sdk.ErrInternal("total debt for this collateral type can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if collateralState.TotalDebt.GT(p.GetDebtParam(changeInDebt.Denom).DebtLimit.AmountOf(changeInDebt.Denom)) {
		return sdk.ErrInternal("change to CSDT would put the system over the debt limit for this debt type")
	}

	k.SetCollateralState(ctx, collateralState)
	return nil
}

// Add/Subtract from collateral debt limit
func (k Keeper) validateAndSetCollateralState(ctx sdk.Context, changeInDebt sdk.Coin) sdk.Error {
	p := k.GetParams(ctx)
	collateralState, found := k.GetCollateralState(ctx, changeInDebt.Denom)
	if !found {
		collateralState = types.CollateralState{Denom: changeInDebt.Denom, TotalDebt: sdk.ZeroInt()} // Already checked that this denom is authorized, so ok to create new CollateralState
	}
	collateralState.TotalDebt = collateralState.TotalDebt.Add(changeInDebt.Amount)
	if collateralState.TotalDebt.IsNegative() {
		return sdk.ErrInternal("total debt for this collateral type can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if collateralState.TotalDebt.GT(p.GetCollateralParam(changeInDebt.Denom).DebtLimit.AmountOf(changeInDebt.Denom)) {
		return sdk.ErrInternal("change to CSDT would put the system over the debt limit for this collateral type")
	}

	k.SetCollateralState(ctx, collateralState)
	return nil
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
