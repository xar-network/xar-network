package denominations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

type GenesisState struct {
	POA string
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func NewGenesisState(poa string) GenesisState {
	return GenesisState{POA: poa}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		POA: "",
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, types.NewParams(""))
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	return GenesisState{POA: params.POA}
}
