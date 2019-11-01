/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// Handle all csdt messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateOrModifyCSDT:
			return handleMsgCreateOrModifyCSDT(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized csdt msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateOrModifyCSDT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateOrModifyCSDT) sdk.Result {

	err := keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, msg.CollateralChange, msg.DebtChange)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
