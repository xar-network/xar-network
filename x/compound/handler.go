package compound

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/compound/internal/keeper"
	"github.com/zar-network/zar-network/x/compound/internal/types"
)

// Handle all compound messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateOrModifyCompound:
			return handleMsgCreateOrModifyCompound(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized compound msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateOrModifyCompound(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateOrModifyCompound) sdk.Result {

	// Can include multiple denominations, but only 1 change per event
	err := keeper.ModifyCompound(ctx, msg.Sender, msg.CollateralDenom, msg.CollateralChange, msg.DebtChange)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
