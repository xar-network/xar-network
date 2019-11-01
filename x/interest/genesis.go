package interest

import (
	"github.com/xar-network/xar-network/x/interest/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	InterestState InterestState `json:"assets" yaml:"assets"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(state InterestState) GenesisState {
	return GenesisState{
		InterestState: state,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		InterestState: DefaultInterestState(),
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetState(ctx, data.InterestState)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	state := keeper.GetState(ctx)
	return NewGenesisState(state)
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := types.ValidateInterestState(data.InterestState)
	if err != nil {
		return err
	}

	return nil
}
