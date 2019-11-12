package issuer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/issuer/internal/keeper"
	"github.com/xar-network/xar-network/x/issuer/internal/types"
)

// TODO Accept Keeper argument
func newHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgIncreaseCredit:
			return handleMsgIncreaseCredit(ctx, msg, k)
		case types.MsgDecreaseCredit:
			return handleMsgDecreaseCredit(ctx, msg, k)
		case types.MsgRevokeLiquidityProvider:
			return handleMsgRevokeLiquidityProvider(ctx, msg, k)
		case types.MsgSetInterest:
			return handleMsgSetInterest(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized issuance Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetInterest(ctx sdk.Context, msg types.MsgSetInterest, k keeper.Keeper) sdk.Result {
	return k.SetInterestRate(ctx, msg.Issuer, msg.InterestRate, msg.Denom)
}

func handleMsgRevokeLiquidityProvider(ctx sdk.Context, msg types.MsgRevokeLiquidityProvider, k keeper.Keeper) sdk.Result {
	return k.RevokeLiquidityProvider(ctx, msg.LiquidityProvider, msg.Issuer)
}

func handleMsgDecreaseCredit(ctx sdk.Context, msg types.MsgDecreaseCredit, k keeper.Keeper) sdk.Result {
	return k.DecreaseCreditOfLiquidityProvider(ctx, msg.LiquidityProvider, msg.Issuer, msg.CreditDecrease)
}

func handleMsgIncreaseCredit(ctx sdk.Context, msg types.MsgIncreaseCredit, k keeper.Keeper) sdk.Result {
	return k.IncreaseCreditOfLiquidityProvider(ctx, msg.LiquidityProvider, msg.Issuer, msg.CreditIncrease)
}
