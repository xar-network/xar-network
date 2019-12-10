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

package synthetic

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/synthetic/internal/keeper"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
)

// Handle all csdt messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgBuySynthetic:
			return handleMsgBuySynthetic(ctx, keeper, msg)
		case types.MsgSellSynthetic:
			return handleMsgSellSynthetic(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized csdt msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgBuySynthetic(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuySynthetic) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.BuySynthetic(ctx, msg.Sender, msg.Coin)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSellSynthetic(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSellSynthetic) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.SellSynthetic(ctx, msg.Sender, msg.Coin)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
