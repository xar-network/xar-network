package denominations

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "assetmanagement" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssueToken:
			return handleMsgIssueToken(ctx, keeper, msg)
		case MsgMintCoins:
			return handleMsgMintCoins(ctx, keeper, msg)
		case MsgBurnCoins:
			return handleMsgBurnCoins(ctx, keeper, msg)
		case MsgFreezeCoins:
			return handleMsgFreezeCoins(ctx, keeper, msg)
		case MsgUnfreezeCoins:
			return handleMsgUnfreezeCoins(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized assetmanagement Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handle message to issue token
func handleMsgIssueToken(ctx sdk.Context, keeper Keeper, msg MsgIssueToken) sdk.Result {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("had to recover when issuing new token: %v", r)
			ctx.Logger().Error(err)
		}
	}()

	// must be lowercase otherwise NewToken will panic
	var newSymbol = strings.ToLower(msg.Symbol)

	token := NewToken(msg.Name, newSymbol, msg.OriginalSymbol, msg.TotalSupply, msg.SourceAddress, msg.Mintable)

	keeperErr := keeper.CoinKeeper.SetCoins(ctx, msg.SourceAddress, token.TotalSupply)
	if keeperErr != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to store new token in bank: %s", keeperErr)).Result()
	}

	err := keeper.SetToken(ctx, newSymbol, token)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to store new token: '%s'", err)).Result()
	}

	newSymbolLog := fmt.Sprintf("new_symbol=%s", newSymbol)
	ctx.Logger().Info(newSymbolLog)
	return sdk.Result{
		Log: newSymbolLog,
	}
}

// handle message to mint coins
func handleMsgMintCoins(ctx sdk.Context, keeper Keeper, msg MsgMintCoins) sdk.Result {
	owner, err := keeper.GetOwner(ctx, msg.Symbol)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", msg.Symbol)).Result()
	}
	if !msg.Owner.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins, err := keeper.CoinKeeper.AddCoins(ctx, owner,
		sdk.NewCoins(sdk.NewInt64Coin(msg.Symbol, msg.Amount)))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to mint coins: '%s'", err)).Result()
	}

	err = keeper.SetTotalSupply(ctx, msg.Symbol, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to set total supply when minting coins: '%s'", err)).Result()
	}
	return sdk.Result{}
}

// handle message to burn coins
func handleMsgBurnCoins(ctx sdk.Context, keeper Keeper, msg MsgBurnCoins) sdk.Result {
	owner, err := keeper.GetOwner(ctx, msg.Symbol)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", msg.Symbol)).Result()
	}
	if !msg.Owner.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins, err := keeper.CoinKeeper.SubtractCoins(ctx, owner,
		sdk.NewCoins(sdk.NewInt64Coin(msg.Symbol, msg.Amount)))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to burn coins: '%s'", err)).Result()
	}

	err = keeper.SetTotalSupply(ctx, msg.Symbol, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to set total supply when burning coins: '%s'", err)).Result()
	}
	return sdk.Result{}
}

// handle message to freeze coins for specific wallet
func handleMsgFreezeCoins(ctx sdk.Context, keeper Keeper, msg MsgFreezeCoins) sdk.Result {
	owner := msg.Owner

	// Todo: Validate you are allowed access to account?
	var customAccount, ok = keeper.AccountKeeper.GetAccount(ctx, owner).(CustomAccount)
	if !ok {
		return sdk.ErrInternal("failed to get correct account type to freeze coins").Result()
	}
	err := customAccount.FreezeCoins(sdk.Coins{sdk.NewInt64Coin(msg.Symbol, msg.Amount)})
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to freeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	keeper.AccountKeeper.SetAccount(ctx, customAccount)
	return sdk.Result{}
}

// handle message to freeze coins for specific wallet
func handleMsgUnfreezeCoins(ctx sdk.Context, keeper Keeper, msg MsgUnfreezeCoins) sdk.Result {
	owner := msg.Owner

	// Todo: Validate you are allowed access to account?
	var customAccount, ok = keeper.AccountKeeper.GetAccount(ctx, owner).(CustomAccount)
	if !ok {
		return sdk.ErrInternal("failed to get correct account type to unfreeze").Result()
	}
	err := customAccount.UnfreezeCoins(sdk.Coins{sdk.NewInt64Coin(msg.Symbol, msg.Amount)})
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to unfreeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	keeper.AccountKeeper.SetAccount(ctx, customAccount)
	return sdk.Result{}
}
