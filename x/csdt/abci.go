package csdt

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
)

// BeginBlocker accrues interest for distribution during EndBlock
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousBlock := k.GetLastAccrualBlock(ctx)
		k.AccrueInterest(ctx, previousBlock)
	}
}
