package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

type Keeper struct {
	ak            auth.AccountKeeper
	sk            supply.Keeper
	paramSubspace subspace.Subspace
	codespace     sdk.CodespaceType
}

func NewKeeper(ak auth.AccountKeeper, sk supply.Keeper, paramstore subspace.Subspace, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		ak:            ak,
		sk:            sk,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		codespace:     codespace,
	}
}

func (k Keeper) Burn(ctx sdk.Context, msg types.MsgBurn) sdk.Result {
	if !k.IsPOA(ctx, msg.From) {
		panic(types.ErrNoPOA(k.codespace))
	} else {
		if !msg.Coins.IsValid() {
			return sdk.ErrInvalidCoins(msg.Coins.String()).Result()
		}
		acc := k.ak.GetAccount(ctx, msg.From)
		if acc == nil {
			return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", msg.From.String())).Result()
		}

		_, hasNeg := acc.GetCoins().SafeSub(msg.Coins)
		if hasNeg {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", acc.GetCoins(), msg.Coins)).Result()
		}

		err := k.sk.SendCoinsFromAccountToModule(ctx, msg.From, types.ModuleName, msg.Coins)
		if err != nil {
			return err.Result()
		}

		err = k.sk.BurnCoins(ctx, types.ModuleName, msg.Coins)
		if err != nil {
			return err.Result()
		}

		return sdk.Result{Events: ctx.EventManager().Events()}
	}
}

func (k Keeper) Mint(ctx sdk.Context, msg types.MsgMint) sdk.Result {
	if !k.IsPOA(ctx, msg.From) {
		panic(types.ErrNoPOA(k.codespace))
	} else {
		if !msg.Coins.IsValid() {
			return sdk.ErrInvalidCoins(msg.Coins.String()).Result()
		}
		acc := k.ak.GetAccount(ctx, msg.From)
		if acc == nil {
			return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", msg.From.String())).Result()
		}

		err := k.sk.MintCoins(ctx, types.ModuleName, msg.Coins)
		if err != nil {
			return err.Result()
		}

		err = k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.From, msg.Coins)
		if err != nil {
			return err.Result()
		}

		return sdk.Result{Events: ctx.EventManager().Events()}
	}
}

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) IsPOA(ctx sdk.Context, from sdk.AccAddress) bool {
	params := k.GetParams(ctx)
	return params.POA == from.String()
}
