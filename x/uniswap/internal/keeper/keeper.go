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
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

// Keeper of the coinswap store
// TODO: make more abstract. replace all keeper references (sk,ak) to related interfaces
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         bank.Keeper
	sk         supply.Keeper
	ak         *auth.AccountKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a uniswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk bank.Keeper, sk supply.Keeper, ak *auth.AccountKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		sk:         sk,
		ak:         ak,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// CreateReservePool initializes a new reserve pool by creating a
// ModuleAccount with minting and burning permissions.
//func (keeper Keeper) CreateReservePool(ctx sdk.Context, moduleName string) {
//	moduleAcc := keeper.sk.GetModuleAccount(ctx, moduleName)
//	if moduleAcc != nil {
//		panic(fmt.Sprintf("reserve pool for %s already exists", moduleName))
//	}
//
//	if _, found := keeper.GetReservePool(ctx, moduleName); found {
//		panic(fmt.Sprintf("reserve pool for %s already exists", moduleName))
//	}
//
//	moduleAcc = supply.NewEmptyModuleAccount(moduleName, supply.Minter, supply.Burner)
//	keeper.sk.SetModuleAccount(ctx, moduleAcc)
//}

// creates new reserve pool and verifies it was created successfully
//func newReservePool(ctx sdk.Context, moduleName string, keeper Keeper) sdk.Error {
//	keeper.CreateReservePool(ctx, moduleName)
//
//	if _, found := keeper.GetReservePool(ctx, moduleName); !found {
//		return types.ErrCannotCreateReservePool(types.DefaultCodespace)
//	}
//
//	return nil
//}

// HasCoins returns whether or not an account has at least coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins ...sdk.Coin) bool {
	return keeper.bk.HasCoins(ctx, addr, coins)
}

// BurnCoins burns liquidity coins from the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (keeper Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Int) sdk.Error {
	mAcc := keeper.ModuleAccountFromName(ctx, moduleName)
	if !mAcc.HasPermission(supply.Burner) {
		return types.ErrInvalidAccountPermission(types.DefaultCodespace, types.MsgAccPermissionsError(moduleName))
	}

	coinsToBurn := sdk.Coins{sdk.Coin{Denom: mAcc.Name, Amount: amt}}
	coinsToBurn.Sort()
	_, err := keeper.bk.SubtractCoins(ctx, mAcc.GetAddress(), coinsToBurn)
	if err != nil {
		return err
	}

	supp := keeper.sk.GetSupply(ctx)
	supp = supp.Deflate(coinsToBurn)
	keeper.sk.SetSupply(ctx, supp)

	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))
	return nil
}

func (keeper Keeper) DoubleSwap(ctx sdk.Context, msg types.MsgSwapOrder) sdk.Result {
	nativeDenom := keeper.GetNativeDenom(ctx)

	if msg.IsBuyOrder {
		nativeMideatorAmt, nonNativeCoinAmt := keeper.DoubleSwapOutputAmount(ctx, msg.Input, msg.Output)
		nativeMideator := sdk.NewCoin(nativeDenom, nativeMideatorAmt)

		moduleNameA := keeper.MustGetPoolName(nativeDenom, msg.Input.Denom)
		mAccA := keeper.ModuleAccountFromName(ctx, moduleNameA)

		moduleNameB := keeper.MustGetPoolName(nativeDenom, msg.Output.Denom)
		mAccB := keeper.ModuleAccountFromName(ctx, moduleNameB)

		err := keeper.SendCoins(ctx, msg.Sender, mAccA.Address, msg.Input)
		if err != nil {
			return err.Result()
		}

		err = keeper.SendCoins(ctx, mAccA.Address, mAccB.Address, nativeMideator)
		if err != nil {
			return err.Result()
		}

		err = keeper.SendCoins(ctx, mAccB.Address, msg.Recipient, sdk.NewCoin(msg.Output.Denom, nonNativeCoinAmt))
		if err != nil {
			return err.Result()
		}
	} else {
		//inputAmountA, inputAmountB := keeper.DoubleSwapOutputAmount(ctx, msg.Input, msg.Output)
	}

	return sdk.Result{}
}

// MintCoins mints liquidity coins to the address and returns liquidity coins
func (keeper Keeper) MintCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Coins, sdk.Error) {
	mAcc := keeper.ak.GetAccount(ctx, addr)
	moduleAcc, ok := mAcc.(*supply.ModuleAccount)
	if !ok {
		return nil, types.ErrInvalidAccountAddr(types.DefaultCodespace, "")
	}

	if !moduleAcc.HasPermission(supply.Minter) {
		msg := fmt.Sprintf("module account %s does not have permissions to mint tokens", moduleAcc.Name)
		return nil, types.ErrInvalidAccountPermission(types.DefaultCodespace, msg)
	}

	coinsToMint := sdk.Coins{sdk.Coin{Denom: moduleAcc.Name, Amount: amt}}
	coinsToMint.Sort()
	_, err := keeper.bk.AddCoins(ctx, mAcc.GetAddress(), coinsToMint)
	if err != nil {
		return nil, err
	}

	//supp := keeper.sk.GetSupply(ctx)
	//supp = supp.Inflate(coinsToMint)
	//
	//keeper.sk.SetSupply(ctx, supp)
	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleAcc.Name))
	return coinsToMint, nil
}

// SendCoin sends coins from the address to the ModuleAccount at moduleName.
func (keeper Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, coins ...sdk.Coin) sdk.Error {
	coinsSorted := sdk.Coins(coins).Sort()
	return keeper.bk.SendCoins(ctx, fromAddr, toAddr, coinsSorted)
}

func (keeper Keeper) SendFromAccToModule(ctx sdk.Context, account sdk.AccAddress, coins sdk.Coins) sdk.Error {
	return keeper.sk.SendCoinsFromAccountToModule(ctx, account, types.ModuleName, coins)
}

func (keeper Keeper) SendFromModuleToAcc(ctx sdk.Context, account sdk.AccAddress, coins sdk.Coins) sdk.Error {
	return keeper.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account, coins)
}

func (keeper Keeper) TransferSwappedCoins(ctx sdk.Context, sender, recipient sdk.AccAddress, userCoin sdk.Coin, moduleCoin sdk.Coin) sdk.Error {
	err := keeper.sk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{userCoin})
	if err != nil {
		return err
	}

	err = keeper.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.Coins{moduleCoin})
	return err
}

func (keeper Keeper) AddLiquidityTransfer(ctx sdk.Context, account sdk.AccAddress, coins sdk.Coins, liquidityVouchers sdk.Coin) sdk.Error {
	err := keeper.SendFromAccToModule(ctx, account, coins)
	if err != nil {
		return err
	}

	err = keeper.sk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(liquidityVouchers))
	if err != nil {
		return err
	}

	err = keeper.SendFromModuleToAcc(ctx, account, sdk.NewCoins(liquidityVouchers))
	return err
}

func (keeper Keeper) RemoveLiquidityTransfer(ctx sdk.Context, account sdk.AccAddress, coins sdk.Coins, userVouchers sdk.Coin) sdk.Error {
	vouchers := sdk.NewCoins(userVouchers)
	totalVouchersAmt := userVouchers.Amount.Mul(sdk.NewInt(2))
	totalVouchers := sdk.NewCoins(sdk.NewCoin(userVouchers.Denom, totalVouchersAmt))
	err := keeper.SendFromAccToModule(ctx, account, vouchers)
	if err != nil {
		return err
	}

	err = keeper.sk.BurnCoins(ctx, types.ModuleName, totalVouchers)
	if err != nil {
		return err
	}

	err = keeper.SendFromModuleToAcc(ctx, account, coins)
	return err
}

// RecieveCoin sends coins from the ModuleAccount at moduleName to the
// address provided.
func (keeper Keeper) RecieveCoins(ctx sdk.Context, addr sdk.AccAddress, coins ...sdk.Coin) sdk.Error {
	// the logic below was probably incorrect too
	// following uniswap whitepaper (https://hackmd.io/@477aQ9OrQTCbVR3fq1Qzxg/HJ9jLsfTz?type=view#Adding-Liquidity)
	// minted tokens are added to both reservePool coins and liquidity provider storage

	//err := keeper.sk.SendCoinsFromModuleToAccount(ctx, moduleName, addr, coins)
	//if err != nil {
	//	panic(err)
	//}

	_, err := keeper.bk.AddCoins(ctx, addr, coins)
	if err != nil {
		return err
	}

	return nil
}

// getReservePoolFromSk returns the total balance of an reserve pool at the
// provided denomination.
func (keeper Keeper) getReservePoolFromSk(ctx sdk.Context, moduleName string) (coins sdk.Coins, found bool) {
	acc := keeper.sk.GetModuleAccount(ctx, moduleName)
	if acc != nil {
		return acc.GetCoins(), true
	}
	return coins, false
}

func (keeper Keeper) GetReservePoolFromAk(ctx sdk.Context, moduleName string) (coins sdk.Coins, found bool) {
	accounts := keeper.ak.GetAllAccounts(ctx)
	for _, v := range accounts {
		acc, ok := v.(*supply.ModuleAccount)
		if !ok {
			continue
		}
		if acc.Name == moduleName {
			return acc.Coins, true
		}
	}
	return coins, false
}

func (keeper Keeper) ModuleAccountFromName(ctx sdk.Context, moduleName string) *supply.ModuleAccount {
	var requestedAcc *supply.ModuleAccount
	fn := func(account exported.Account) (stop bool) {
		acc, ok := account.(*supply.ModuleAccount)
		if !ok {
			return false
		}
		if acc.Name != moduleName {
			return false
		}

		requestedAcc = acc
		return true
	}

	keeper.ak.IterateAccounts(ctx, fn)
	return requestedAcc
}

//func (keeper Keeper) GetReservePool(ctx sdk.Context, moduleName string) (sdk.Coins, bool) {
//	rp, found := keeper.getReservePoolFromSk(ctx, moduleName)
//	if found {
//		return rp, found
//	}
//
//	return keeper.GetReservePoolFromAk(ctx, moduleName)
//}

func (keeper Keeper) AddInitialLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity) sdk.Result {
	nativeDenom, _, moduleName := keeper.MustGetAllDenoms(ctx, msg)

	coinAmount := msg.Deposit.Amount.BigInt()
	nativeAmount := msg.DepositAmount.BigInt()
	mintAmtBigint := (coinAmount.Mul(coinAmount, nativeAmount)).Sqrt(coinAmount)
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	amtToMint := sdk.NewIntFromBigInt(mintAmtBigint)

	if !keeper.HasCoins(ctx, msg.Sender, nativeCoinDeposited, msg.Deposit) {
		return sdk.ErrInsufficientCoins(types.InsufficientCoins).Result()
	}

	mAcc := keeper.ModuleAccountFromName(ctx, moduleName)
	err := keeper.transferLiquidityCoins(ctx, msg, msg.Deposit, nativeCoinDeposited, amtToMint, mAcc)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// creates new reserve pool and verifies it was created successfully
//func (keeper Keeper) NewReservePool(ctx sdk.Context, moduleName string) sdk.Error {
//	keeper.CreateReservePool(ctx, moduleName)
//
//	if _, found := keeper.GetReservePool(ctx, moduleName); !found {
//		return types.ErrCannotCreateReservePool(types.DefaultCodespace)
//	}
//
//	return nil
//}

func (keeper Keeper) AddLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity, reservePool sdk.Coins) sdk.Result {
	nativeDenom, _, moduleName := keeper.MustGetAllDenoms(ctx, msg)
	nativeBalance := reservePool.AmountOf(nativeDenom)
	liquidityCoinBalance := reservePool.AmountOf(moduleName)
	if liquidityCoinBalance.LTE(sdk.NewInt(0)) {
		return types.ErrInsufficientLiquidityAmount(types.DefaultCodespace).Result()
	}

	amtToMint := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance)
	coinAmountDeposited := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance)
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	coinDeposited := sdk.NewCoin(msg.Deposit.Denom, coinAmountDeposited)

	if !keeper.HasCoins(ctx, msg.Sender, nativeCoinDeposited, coinDeposited) {
		return sdk.ErrInsufficientCoins(types.InsufficientCoins).Result()
	}

	mAcc := keeper.ModuleAccountFromName(ctx, moduleName)
	err := keeper.transferLiquidityCoins(ctx, msg, nativeCoinDeposited, coinDeposited, amtToMint, mAcc)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func (keeper Keeper) transferLiquidityCoins(ctx sdk.Context, msg *types.MsgAddLiquidity, nativeCoin, coin sdk.Coin, amtToMint sdk.Int, moduleAcc *supply.ModuleAccount) sdk.Error {
	err := keeper.SendCoins(ctx, msg.Sender, moduleAcc.Address, nativeCoin, coin)
	if err != nil {
		return err
	}

	// mint liquidity vouchers for sender
	mintCoins, err := keeper.MintCoins(ctx, moduleAcc.Address, amtToMint)
	if err != nil {
		return err
	}

	return keeper.RecieveCoins(ctx, msg.Sender, mintCoins...)
}

// GetNativeDenom returns the native denomination for this module from the
// global param store.
func (keeper Keeper) GetNativeDenom(ctx sdk.Context) (nativeDenom string) {
	keeper.paramSpace.Get(ctx, types.KeyNativeDenom, &nativeDenom)
	return
}

func (keeper Keeper) MustGetAllDenoms(ctx sdk.Context, msg *types.MsgAddLiquidity) (nativeDenom string, tokenDenom string, moduleName string) {
	nativeDenom = keeper.GetNativeDenom(ctx)
	tokenDenom = msg.Deposit.Denom

	return nativeDenom, tokenDenom, keeper.MustGetPoolName(nativeDenom, tokenDenom)
}

// GetFeeParam returns the current FeeParam from the global param store
func (keeper Keeper) GetFeeParam(ctx sdk.Context) (feeParam types.FeeParam) {
	keeper.paramSpace.Get(ctx, types.KeyFee, &feeParam)
	return
}

// SetParams sets the parameters for the coinswap module.
func (keeper Keeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
