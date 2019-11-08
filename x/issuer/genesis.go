package issuer

import (
	"github.com/xar-network/xar-network/x/issuer/internal/keeper"
	"github.com/xar-network/xar-network/x/issuer/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type genesisState struct {
	Issuers []types.Issuer
}

func initGenesis(ctx sdk.Context, k keeper.Keeper, state genesisState) {
	for _, issuer := range state.Issuers {
		k.AddIssuer(ctx, issuer)
	}
}

func defaultGenesisState() genesisState {
	return genesisState{}
}
