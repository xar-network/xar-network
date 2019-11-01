package keeper

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// Keeper csdt Keeper
type Keeper struct {
	storeKey       sdk.StoreKey
	pricefeed      pricefeedKeeper
	bank           bankKeeper
	paramsSubspace params.Subspace
	cdc            *codec.Codec
}

// NewKeeper creates a new keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, subspace params.Subspace, pricefeed pricefeedKeeper, bank bankKeeper) Keeper {
	subspace = subspace.WithKeyTable(types.CreateParamsKeyTable())
	return Keeper{
		storeKey:       storeKey,
		pricefeed:      pricefeed,
		bank:           bank,
		paramsSubspace: subspace,
		cdc:            cdc,
	}
}

// ModifyCSDT creates, changes, or deletes a CSDT
// TODO can/should this function be split up?
func (k Keeper) ModifyCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, changeInCollateral sdk.Int, changeInDebt sdk.Int) sdk.Error {

	// Phase 1: Get state, make changes in memory and check if they're ok.

	// Check collateral type ok
	p := k.GetParams(ctx)
	if !p.IsCollateralPresent(collateralDenom) { // maybe abstract this logic into GetCSDT
		return sdk.ErrInternal("collateral type not enabled to create CSDTs")
	}

	// Check the owner has enough collateral and stable coins
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

	// Change collateral and debt recorded in CSDT
	// Get CSDT (or create if not exists)
	csdt, found := k.GetCSDT(ctx, owner, collateralDenom)
	if !found {
		csdt = types.CSDT{Owner: owner, CollateralDenom: collateralDenom, CollateralAmount: sdk.ZeroInt(), Debt: sdk.ZeroInt()}
	}
	// Add/Subtract collateral and debt
	csdt.CollateralAmount = csdt.CollateralAmount.Add(changeInCollateral)
	if csdt.CollateralAmount.IsNegative() {
		return sdk.ErrInternal(" can't withdraw more collateral than exists in CSDT")
	}
	csdt.Debt = csdt.Debt.Add(changeInDebt)
	if csdt.Debt.IsNegative() {
		return sdk.ErrInternal("can't pay back more debt than exists in CSDT")
	}
	isUnderCollateralized := csdt.IsUnderCollateralized(
		k.pricefeed.GetCurrentPrice(ctx, csdt.CollateralDenom).Price,
		p.GetCollateralParams(csdt.CollateralDenom).LiquidationRatio,
	)
	if isUnderCollateralized {
		return sdk.ErrInternal("Change to CSDT would put it below liquidation ratio")
	}
	// TODO check for dust

	// Add/Subtract from global debt limit
	gDebt := k.GetGlobalDebt(ctx)
	gDebt = gDebt.Add(changeInDebt)
	if gDebt.IsNegative() {
		return sdk.ErrInternal("global debt can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if gDebt.GT(p.GlobalDebtLimit) {
		return sdk.ErrInternal("change to CSDT would put the system over the global debt limit")
	}

	// Add/Subtract from collateral debt limit
	collateralState, found := k.GetCollateralState(ctx, csdt.CollateralDenom)
	if !found {
		collateralState = types.CollateralState{Denom: csdt.CollateralDenom, TotalDebt: sdk.ZeroInt()} // Already checked that this denom is authorized, so ok to create new CollateralState
	}
	collateralState.TotalDebt = collateralState.TotalDebt.Add(changeInDebt)
	if collateralState.TotalDebt.IsNegative() {
		return sdk.ErrInternal("total debt for this collateral type can't be negative") // This should never happen if debt per CSDT can't be negative
	}
	if collateralState.TotalDebt.GT(p.GetCollateralParams(csdt.CollateralDenom).DebtLimit) {
		return sdk.ErrInternal("change to CSDT would put the system over the debt limit for this collateral type")
	}

	// Phase 2: Update all the state

	// change owner's coins (increase or decrease)
	var err sdk.Error
	if changeInCollateral.IsNegative() {
		_, err = k.bank.AddCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral.Neg())))
	} else {
		_, err = k.bank.SubtractCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(collateralDenom, changeInCollateral)))
	}
	if err != nil {
		panic(err) // this shouldn't happen because coin balance was checked earlier
	}
	if changeInDebt.IsNegative() {
		_, err = k.bank.SubtractCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt.Neg())))
	} else {
		_, err = k.bank.AddCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin(types.StableDenom, changeInDebt)))
	}
	if err != nil {
		panic(err) // this shouldn't happen because coin balance was checked earlier
	}
	// Set CSDT
	if csdt.CollateralAmount.IsZero() && csdt.Debt.IsZero() { // TODO maybe abstract this logic into setCSDT
		k.deleteCSDT(ctx, csdt)
	} else {
		k.setCSDT(ctx, csdt)
	}
	// set total debts
	k.SetGlobalDebt(ctx, gDebt)
	k.setCollateralState(ctx, collateralState)

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
		k.pricefeed.GetCurrentPrice(ctx, csdt.CollateralDenom).Price,
		p.GetCollateralParams(csdt.CollateralDenom).LiquidationRatio,
	)
	if !isUnderCollateralized {
		return sdk.ErrInternal("CSDT is not currently under the liquidation ratio")
	}

	// Remove Collateral
	if collateralToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative collateral")
	}
	csdt.CollateralAmount = csdt.CollateralAmount.Sub(collateralToSeize)
	if csdt.CollateralAmount.IsNegative() {
		return sdk.ErrInternal("can't seize more collateral than exists in CSDT")
	}

	// Remove Debt
	if debtToSeize.IsNegative() {
		return sdk.ErrInternal("cannot seize negative debt")
	}
	csdt.Debt = csdt.Debt.Sub(debtToSeize)
	if csdt.Debt.IsNegative() {
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
	if csdt.CollateralAmount.IsZero() && csdt.Debt.IsZero() { // TODO maybe abstract this logic into setCSDT
		k.deleteCSDT(ctx, csdt)
	} else {
		k.setCSDT(ctx, csdt)
	}
	k.setCollateralState(ctx, collateralState)
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

// ---------- Module Parameters ----------

func (k Keeper) GetParams(ctx sdk.Context) types.CsdtModuleParams {
	var p types.CsdtModuleParams
	k.paramsSubspace.Get(ctx, types.ModuleParamsKey, &p)
	return p
}

// This is only needed to be able to setup the store from the genesis file. The keeper should not change any of the params itself.
func (k Keeper) SetParams(ctx sdk.Context, csdtModuleParams types.CsdtModuleParams) {
	k.paramsSubspace.Set(ctx, types.ModuleParamsKey, &csdtModuleParams)
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
func (k Keeper) setCSDT(ctx sdk.Context, csdt types.CSDT) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(csdt)
	store.Set(k.getCSDTKey(csdt.Owner, csdt.CollateralDenom), bz)
}
func (k Keeper) deleteCSDT(ctx sdk.Context, csdt types.CSDT) { // TODO should this id the csdt by passing in owner,collateralDenom pair?
	// get store
	store := ctx.KVStore(k.storeKey)
	// delete key
	store.Delete(k.getCSDTKey(csdt.Owner, csdt.CollateralDenom))
}

// GetCSDTs returns all CSDTs, optionally filtered by collateral type and liquidation price.
// `price` filters for CSDTs that will be below the liquidation ratio when the collateral is at that specified price.
func (k Keeper) GetCSDTs(ctx sdk.Context, collateralDenom string, price sdk.Dec) (types.CSDTs, sdk.Error) {
	// Validate inputs
	params := k.GetParams(ctx)
	if len(collateralDenom) != 0 && !params.IsCollateralPresent(collateralDenom) {
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
			if csdt.IsUnderCollateralized(price, params.GetCollateralParams(collateralDenom).LiquidationRatio) {
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
func (k Keeper) setCollateralState(ctx sdk.Context, collateralstate types.CollateralState) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(collateralstate)
	store.Set(k.getCollateralStateKey(collateralstate.Denom), bz)
}

// ---------- Weird Bank Stuff ----------
// This only exists because module accounts aren't really a thing yet.
// Also because we need module accounts that allow for burning/minting.

// These functions make the CSDT module act as a bank keeper, ie it fulfills the bank.Keeper interface.
// It intercepts calls to send coins to/from the liquidator module account, otherwise passing the calls onto the normal bank keeper.

// Not sure if module accounts are good, but they make the auction module more general:
// - startAuction would just "mints" coins, relying on calling function to decrement them somewhere
// - closeAuction would have to call something specific for the receiver module to accept coins (like liquidationKeeper.AddStableCoins)

// The auction and liquidator modules can probably just use SendCoins to keep things safe (instead of AddCoins and SubtractCoins).
// So they should define their own interfaces which this module should fulfill, rather than this fulfilling the entire bank.Keeper interface.

// bank.Keeper interfaces:
// type SendKeeper interface {
// 	type ViewKeeper interface {
// 		GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
// 		HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
// 		Codespace() sdk.CodespaceType
// 	}
// 	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
// 	GetSendEnabled(ctx sdk.Context) bool
// 	SetSendEnabled(ctx sdk.Context, enabled bool)
// }
// type Keeper interface {
// 	SendKeeper
// 	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
// 	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
// 	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
// 	InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error)
// 	DelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
// 	UndelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)

var LiquidatorAccountAddress = sdk.AccAddress([]byte("whatever"))
var liquidatorAccountKey = []byte("liquidatorAccount")

func (k Keeper) GetLiquidatorAccountAddress() sdk.AccAddress {
	return LiquidatorAccountAddress
}

type LiquidatorModuleAccount struct {
	Coins sdk.Coins // keeps track of seized collateral, surplus csdt, and mints/burns gov coins
}

func (k Keeper) AddCoins(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) (sdk.Coins, sdk.Error) {
	// intercept module account
	if address.Equals(LiquidatorAccountAddress) {
		if !amount.IsValid() {
			return nil, sdk.ErrInvalidCoins(amount.String())
		}
		// remove gov token from list
		filteredCoins := stripGovCoin(amount)
		// add coins to module account
		lma := k.getLiquidatorModuleAccount(ctx)
		updatedCoins := lma.Coins.Add(filteredCoins)
		if updatedCoins.IsAnyNegative() {
			return amount, sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", lma.Coins, amount))
		}
		lma.Coins = updatedCoins
		k.setLiquidatorModuleAccount(ctx, lma)
		return updatedCoins, nil
	} else {
		return k.bank.AddCoins(ctx, address, amount)
	}
}

// TODO abstract stuff better
func (k Keeper) SubtractCoins(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) (sdk.Coins, sdk.Error) {
	// intercept module account
	if address.Equals(LiquidatorAccountAddress) {
		if !amount.IsValid() {
			return nil, sdk.ErrInvalidCoins(amount.String())
		}
		// remove gov token from list
		filteredCoins := stripGovCoin(amount)
		// subtract coins from module account
		lma := k.getLiquidatorModuleAccount(ctx)
		updatedCoins, isNegative := lma.Coins.SafeSub(filteredCoins)
		if isNegative {
			return amount, sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", lma.Coins, amount))
		}
		lma.Coins = updatedCoins
		k.setLiquidatorModuleAccount(ctx, lma)
		return updatedCoins, nil
	} else {
		return k.bank.SubtractCoins(ctx, address, amount)
	}
}

// TODO Should this return anything for the gov coin balance? Currently returns nothing.
func (k Keeper) GetCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	if address.Equals(LiquidatorAccountAddress) {
		return k.getLiquidatorModuleAccount(ctx).Coins
	} else {
		return k.bank.GetCoins(ctx, address)
	}
}

// TODO test this with unsorted coins
func (k Keeper) HasCoins(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) bool {
	if address.Equals(LiquidatorAccountAddress) {
		return true
	} else {
		return k.getLiquidatorModuleAccount(ctx).Coins.IsAllGTE(stripGovCoin(amount))
	}
}

func (k Keeper) getLiquidatorModuleAccount(ctx sdk.Context) LiquidatorModuleAccount {
	// get store
	store := ctx.KVStore(k.storeKey)
	// get bytes
	bz := store.Get(liquidatorAccountKey)
	if bz == nil {
		return LiquidatorModuleAccount{} // TODO is it safe to do this, or better to initialize the account explicitly
	}
	// unmarshal
	var lma LiquidatorModuleAccount
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lma)
	return lma
}
func (k Keeper) setLiquidatorModuleAccount(ctx sdk.Context, lma LiquidatorModuleAccount) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(lma)
	store.Set(liquidatorAccountKey, bz)
}
func stripGovCoin(coins sdk.Coins) sdk.Coins {
	filteredCoins := sdk.NewCoins()
	for _, c := range coins {
		if c.Denom != types.GovDenom {
			filteredCoins = append(filteredCoins, c)
		}
	}
	return filteredCoins
}
