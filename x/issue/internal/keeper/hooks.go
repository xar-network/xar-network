package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

// Wrapper struct
type Hooks struct {
	keeper Keeper
}

// Create new box hooks
func (keeper Keeper) Hooks() Hooks { return Hooks{keeper} }

func (hooks Hooks) CanSend(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (bool, sdk.Error) {
	for _, v := range amt {
		if !types.IsIssueId(v.Denom) {
			continue
		}
		if err := hooks.keeper.CheckFreeze(ctx, fromAddr, toAddr, v.Denom); err != nil {
			return false, err
		}
	}
	return true, nil
}
