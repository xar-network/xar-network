package authority

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	AuthorityKey sdk.AccAddress `json:"authority_key" yaml:"authority_key"`
}

func NewGenesisState(authorityKey sdk.AccAddress) GenesisState {
	return GenesisState{
		AuthorityKey: authorityKey,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, state GenesisState) {
	keeper.SetAuthority(ctx, state.AuthorityKey)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	authorityKey := keeper.GetAuthority(ctx)
	return NewGenesisState(authorityKey)
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return nil
}
