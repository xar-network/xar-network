/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

package coinswap

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwapOrder:
			return HandleMsgSwapOrder(ctx, msg, k)
		case MsgAddLiquidity:
			return HandleMsgAddLiquidity(ctx, msg, k)
		case MsgRemoveLiquidity:
			return HandleMsgRemoveLiquidity(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized coinswap message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgSwapOrder.
func HandleMsgSwapOrder(ctx sdk.Context, msg MsgSwapOrder, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgSwapOrder").Result()
	}
	tags, err := k.HandleSwap(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Tags: tags}
}

// HandleMsgAddLiquidity If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
func HandleMsgAddLiquidity(ctx sdk.Context, msg MsgAddLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgAddLiquidity").Result()
	}

	tags, err := k.HandleAddLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, msg MsgRemoveLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgRemoveLiquidity").Result()
	}

	tags, err := k.HandleRemoveLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
