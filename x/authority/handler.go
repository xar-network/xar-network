package authority

import (
	"fmt"

	"github.com/xar-network/xar-network/x/authority/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func newHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result sdk.Result) {
		defer func() {
			if r := recover(); r != nil {
				switch o := r.(type) {
				case sdk.Result:
					result = o
				case sdk.Error:
					result = o.Result()
				default:
					panic(r)
				}
			}
		}()

		switch msg := msg.(type) {
		case types.MsgCreateIssuer:
			return k.CreateIssuer(ctx, msg.Authority, msg.Issuer, msg.Denominations)
		case types.MsgCreateOracle:
			return k.CreateOracle(ctx, msg.Authority, msg.Oracle)
		case types.MsgDestroyIssuer:
			return k.DestroyIssuer(ctx, msg.Authority, msg.Issuer)
		default:
			errMsg := fmt.Sprintf("Unrecognized authority Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
