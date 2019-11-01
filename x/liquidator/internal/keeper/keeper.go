package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

type Keeper struct {
	cdc            *codec.Codec
	paramsSubspace params.Subspace
	storeKey       sdk.StoreKey
	csdtKeeper      csdtKeeper
	auctionKeeper  auctionKeeper
	bankKeeper     bankKeeper
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, subspace params.Subspace, csdtKeeper csdtKeeper, auctionKeeper auctionKeeper, bankKeeper bankKeeper) Keeper {
	subspace = subspace.WithKeyTable(types.CreateParamsKeyTable())
	return Keeper{
		cdc:            cdc,
		paramsSubspace: subspace,
		storeKey:       storeKey,
		csdtKeeper:      csdtKeeper,
		auctionKeeper:  auctionKeeper,
		bankKeeper:     bankKeeper,
	}
}

// SeizeAndStartCollateralAuction pulls collateral out of a CSDT and sells it in an auction for stable coin. Excess collateral goes to the original CSDT owner.
// Known as Cat.bite in maker
// result: stable coin is transferred to module account, collateral is transferred from module account to buyer, (and any excess collateral is transferred to original CSDT owner)
func (k Keeper) SeizeAndStartCollateralAuction(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string) (auction.ID, sdk.Error) {
	// Get CSDT
	csdt, found := k.csdtKeeper.GetCSDT(ctx, owner, collateralDenom)
	if !found {
		return 0, sdk.ErrInternal("CSDT not found")
	}

	// Calculate amount of collateral to sell in this auction
	params := k.GetParams(ctx).GetCollateralParams(csdt.CollateralDenom)
	collateralToSell := sdk.MinInt(csdt.CollateralAmount, params.AuctionSize)
	// Calculate the corresponding maximum amount of stable coin to raise TODO test maths
	stableToRaise := sdk.NewDecFromInt(collateralToSell).Quo(sdk.NewDecFromInt(csdt.CollateralAmount)).Mul(sdk.NewDecFromInt(csdt.Debt)).RoundInt()

	// Seize the collateral and debt from the CSDT
	err := k.partialSeizeCSDT(ctx, owner, collateralDenom, collateralToSell, stableToRaise)
	if err != nil {
		return 0, err
	}

	// Start "forward reverse" auction type
	lot := sdk.NewCoin(csdt.CollateralDenom, collateralToSell)
	maxBid := sdk.NewCoin(k.csdtKeeper.GetStableDenom(), stableToRaise)
	auctionID, err := k.auctionKeeper.StartForwardReverseAuction(ctx, k.csdtKeeper.GetLiquidatorAccountAddress(), lot, maxBid, owner)
	if err != nil {
		panic(err) // TODO how can errors here be handled to be safe with the state update in PartialSeizeCSDT?
	}
	return auctionID, nil
}

// StartDebtAuction sells off minted gov coin to raise set amounts of stable coin.
// Known as Vow.flop in maker
// result: minted gov coin moved to highest bidder, stable coin moved to moduleAccount
func (k Keeper) StartDebtAuction(ctx sdk.Context) (auction.ID, sdk.Error) {

	// Ensure amount of seized stable coin is 0 (ie Joy = 0)
	stableCoins := k.bankKeeper.GetCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress()).AmountOf(k.csdtKeeper.GetStableDenom())
	if !stableCoins.IsZero() {
		return 0, sdk.ErrInternal("debt auction cannot be started as there is outstanding stable coins")
	}

	// check the seized debt is above a threshold
	params := k.GetParams(ctx)
	seizedDebt := k.GetSeizedDebt(ctx)
	if seizedDebt.Available().LT(params.DebtAuctionSize) {
		return 0, sdk.ErrInternal("not enough seized debt to start an auction")
	}
	// start reverse auction, selling minted gov coin for stable coin
	auctionID, err := k.auctionKeeper.StartReverseAuction(
		ctx,
		k.csdtKeeper.GetLiquidatorAccountAddress(),
		sdk.NewCoin(k.csdtKeeper.GetStableDenom(), params.DebtAuctionSize),
		sdk.NewInt64Coin(k.csdtKeeper.GetGovDenom(), 2^255-1), // TODO is there a way to avoid potentially minting infinite gov coin?
	)
	if err != nil {
		return 0, err
	}
	// Record amount of debt sent for auction. Debt can only be reduced in lock step with reducing stable coin
	seizedDebt.SentToAuction = seizedDebt.SentToAuction.Add(params.DebtAuctionSize)
	k.setSeizedDebt(ctx, seizedDebt)
	return auctionID, nil
}

// With no stability and liquidation fees, surplus auctions can never be run.
// StartSurplusAuction sells off excess stable coin in exchange for gov coin, which is burned
// Known as Vow.flap in maker
// result: stable coin removed from module account (eventually to buyer), gov coin transferred to module account
// func (k Keeper) StartSurplusAuction(ctx sdk.Context) (auction.ID, sdk.Error) {

// 	// TODO ensure seized debt is 0

// 	// check there is enough surplus to be sold
// 	surplus := k.bankKeeper.GetCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress()).AmountOf(k.csdtKeeper.GetStableDenom())
// 	if surplus.LT(SurplusAuctionSize) {
// 		return 0, sdk.ErrInternal("not enough surplus stable coin to start an auction")
// 	}
// 	// start normal auction, selling stable coin
// 	auctionID, err := k.auctionKeeper.StartForwardAuction(
// 		ctx,
// 		k.csdtKeeper.GetLiquidatorAccountAddress(),
// 		sdk.NewCoin(k.csdtKeeper.GetStableDenom(), SurplusAuctionSize),
// 		sdk.NewInt64Coin(k.csdtKeeper.GetGovDenom(), 0),
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// Starting the auction will remove coins from the account, so they don't need modified here.
// 	return auctionID, nil
// }

// PartialSeizeCSDT seizes some collateral and debt from an under-collateralized CSDT.
func (k Keeper) partialSeizeCSDT(ctx sdk.Context, owner sdk.AccAddress, collateralDenom string, collateralToSeize sdk.Int, debtToSeize sdk.Int) sdk.Error { // aka Cat.bite
	// Seize debt and collateral in the csdt module. This also validates the inputs.
	err := k.csdtKeeper.PartialSeizeCSDT(ctx, owner, collateralDenom, collateralToSeize, debtToSeize)
	if err != nil {
		return err // csdt could be not found, or not under collateralized, or inputs invalid
	}

	// increment the total seized debt (Awe) by csdt.debt
	seizedDebt := k.GetSeizedDebt(ctx)
	seizedDebt.Total = seizedDebt.Total.Add(debtToSeize)
	k.setSeizedDebt(ctx, seizedDebt)

	// add csdt.collateral amount of coins to the moduleAccount (so they can be transferred to the auction later)
	coins := sdk.NewCoins(sdk.NewCoin(collateralDenom, collateralToSeize))
	_, err = k.bankKeeper.AddCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress(), coins)
	if err != nil {
		panic(err) // TODO this shouldn't happen?
	}
	return nil
}

// SettleDebt removes equal amounts of debt and stable coin from the liquidator's reserves (and also updates the global debt in the csdt module).
// This is called in the handler when a debt or surplus auction is started
// TODO Should this be called with an amount, rather than annihilating the maximum?
func (k Keeper) SettleDebt(ctx sdk.Context) sdk.Error {
	// Calculate max amount of debt and stable coins that can be settled (ie annihilated)
	debt := k.GetSeizedDebt(ctx)
	stableCoins := k.bankKeeper.GetCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress()).AmountOf(k.csdtKeeper.GetStableDenom())
	settleAmount := sdk.MinInt(debt.Total, stableCoins)

	// Call csdt module to reduce GlobalDebt. This can fail if genesis not set
	err := k.csdtKeeper.ReduceGlobalDebt(ctx, settleAmount)
	if err != nil {
		return err
	}

	// Decrement total seized debt (also decrement from SentToAuction debt)
	updatedDebt, err := debt.Settle(settleAmount)
	if err != nil {
		return err // this should not error in this context
	}
	k.setSeizedDebt(ctx, updatedDebt)

	// Subtract stable coin from moduleAccout
	k.bankKeeper.SubtractCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress(), sdk.Coins{sdk.NewCoin(k.csdtKeeper.GetStableDenom(), settleAmount)})
	return nil
}

// ---------- Module Parameters ----------

func (k Keeper) GetParams(ctx sdk.Context) types.LiquidatorModuleParams {
	var params types.LiquidatorModuleParams
	k.paramsSubspace.Get(ctx, types.ModuleParamsKey, &params)
	return params
}

// This is only needed to be able to setup the store from the genesis file. The keeper should not change any of the params itself.
func (k Keeper) SetParams(ctx sdk.Context, params types.LiquidatorModuleParams) {
	k.paramsSubspace.Set(ctx, types.ModuleParamsKey, &params)
}

// ---------- Store Wrappers ----------

func (k Keeper) getSeizedDebtKey() []byte {
	return []byte("seizedDebt")
}
func (k Keeper) GetSeizedDebt(ctx sdk.Context) types.SeizedDebt {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(k.getSeizedDebtKey())
	if bz == nil {
		// TODO make initial seized debt and CSDTs configurable at genesis, then panic here if not found
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(types.SeizedDebt{sdk.ZeroInt(), sdk.ZeroInt()})
	}
	var seizedDebt types.SeizedDebt
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &seizedDebt)
	return seizedDebt
}
func (k Keeper) setSeizedDebt(ctx sdk.Context, debt types.SeizedDebt) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(debt)
	store.Set(k.getSeizedDebtKey(), bz)
}
