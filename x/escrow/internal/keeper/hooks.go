package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/escrow/internal/types"
)

// Wrapper struct
type Hooks struct {
	keeper Keeper
}

// Create new box hooks
func (keeper Keeper) Hooks() Hooks { return Hooks{keeper} }

func (hooks Hooks) CanSend(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (bool, sdk.Error) {
	for _, v := range amt {
		if !types.IsId(v.Denom) {
			continue
		}
		box := hooks.keeper.GetBox(ctx, v.Denom)
		if box == nil {
			continue
		}
		if box.IsTransferDisabled() {
			return false, types.ErrCanNotTransfer(v.Denom)
		}
	}
	return true, nil
}
