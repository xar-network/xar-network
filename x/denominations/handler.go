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

package denominations

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

// NewHandler returns a handler for "assetmanagement" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgIssueToken:
			return handleMsgIssueToken(ctx, keeper, msg)
		case types.MsgMintCoins:
			return handleMsgMintCoins(ctx, keeper, msg)
		case types.MsgBurnCoins:
			return handleMsgBurnCoins(ctx, keeper, msg)
		case types.MsgFreezeCoins:
			return handleMsgFreezeCoins(ctx, keeper, msg)
		case types.MsgUnfreezeCoins:
			return handleMsgUnfreezeCoins(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized assetmanagement Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handle message to issue token
func handleMsgIssueToken(ctx sdk.Context, k Keeper, msg types.MsgIssueToken) sdk.Result {

	// must be lowercase otherwise NewToken will panic
	var newSymbol = strings.ToLower(msg.Symbol)

	token := types.NewToken(
		msg.Name, newSymbol,
		msg.OriginalSymbol,
		msg.MaxSupply,
		msg.Owner,
		msg.Mintable,
	)

	return k.IssueToken(ctx, msg.SourceAddress, msg.Owner, *token)
}

// handle message to mint coins
func handleMsgMintCoins(ctx sdk.Context, keeper Keeper, msg types.MsgMintCoins) sdk.Result {
	return keeper.MintCoins(ctx, msg.Owner, msg.Amount, msg.Symbol)
}

// handle message to burn coins
func handleMsgBurnCoins(ctx sdk.Context, keeper Keeper, msg types.MsgBurnCoins) sdk.Result {
	return keeper.BurnCoins(ctx, msg.Owner, msg.Amount, msg.Symbol)
}

// handle message to freeze coins for specific wallet
func handleMsgFreezeCoins(ctx sdk.Context, keeper Keeper, msg types.MsgFreezeCoins) sdk.Result {
	return keeper.FreezeCoins(ctx, msg.Owner, msg.Address, msg.Amount, msg.Symbol)
}

// handle message to freeze coins for specific wallet
func handleMsgUnfreezeCoins(ctx sdk.Context, keeper Keeper, msg types.MsgUnfreezeCoins) sdk.Result {
	return keeper.UnfreezeCoins(ctx, msg.Owner, msg.Address, msg.Amount, msg.Symbol)
}
