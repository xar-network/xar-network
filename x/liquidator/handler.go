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

package liquidator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidator/internal/keeper"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// Handle all liquidator messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSeizeAndStartCollateralAuction:
			return handleMsgSeizeAndStartCollateralAuction(ctx, keeper, msg)
		case types.MsgStartDebtAuction:
			return handleMsgStartDebtAuction(ctx, keeper)
		// case MsgStartSurplusAuction:
		// 	return handleMsgStartSurplusAuction(ctx, keeper)
		default:
			errMsg := fmt.Sprintf("Unrecognized liquidator msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSeizeAndStartCollateralAuction(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSeizeAndStartCollateralAuction) sdk.Result {
	_, err := keeper.SeizeAndStartCollateralAuction(ctx, msg.CsdtOwner, msg.CollateralDenom)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{} // TODO tags, return auction ID
}

func handleMsgStartDebtAuction(ctx sdk.Context, keeper keeper.Keeper) sdk.Result {
	// cancel out any debt and stable coins before trying to start auction
	keeper.SettleDebt(ctx)
	// start an auction
	_, err := keeper.StartDebtAuction(ctx)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{} // TODO tags, return auction ID
}

// With no stability and liquidation fees, surplus auctions can never be run.
// func handleMsgStartSurplusAuction(ctx sdk.Context, keeper Keeper) sdk.Result {
// 	// cancel out any debt and stable coins before trying to start auction
//  keeper.settleDebt(ctx)
// 	_, err := keeper.StartSurplusAuction(ctx)
// 	if err != nil {
// 		return err.Result()
// 	}
// 	return sdk.Result{} // TODO tags
// }
