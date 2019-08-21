package issue

import (
	"bytes"
	"fmt"

	"github.com/zar-network/zar-network/x/issue/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zar-network/zar-network/x/issue/internal/keeper"
)

// GenesisState - all issue state that must be provided at genesis
type GenesisState struct {
	StartingIssueId uint64          `json:"starting_issue_id"`
	Issues          []CoinIssueInfo `json:"issues"`
	Params          types.Params    `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(startingIssueId uint64) GenesisState {
	params := types.DefaultParams("ftm")
	return GenesisState{
		StartingIssueId: startingIssueId,
		Params:          params,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.CoinIssueMinId)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	err := k.SetInitialIssueStartingIssueId(ctx, data.StartingIssueId)
	if err != nil {
		panic(err)
	}

	err = k.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}

	for _, issue := range data.Issues {
		k.AddIssue(ctx, &issue)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	genesisState := GenesisState{}

	startingIssueId, _ := keeper.PeekCurrentIssueID(ctx)
	genesisState.StartingIssueId = startingIssueId

	genesisState.Params = keeper.GetParams(ctx)
	genesisState.Issues = keeper.ListAll(ctx)

	return genesisState

}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.IssueFee.IsNegative() {
		return fmt.Errorf("invalid issue fee: %s", data.Params.IssueFee.String())
	}
	if data.Params.MintFee.IsNegative() {
		return fmt.Errorf("invalid mint fee: %s", data.Params.MintFee.String())
	}
	if data.Params.BurnFee.IsNegative() {
		return fmt.Errorf("invalid burn fee: %s", data.Params.BurnFee.String())
	}
	if data.Params.BurnFromFee.IsNegative() {
		return fmt.Errorf("invalid burn from fee: %s", data.Params.BurnFromFee.String())
	}
	if data.Params.FreezeFee.IsNegative() {
		return fmt.Errorf("invalid freeze fee: %s", data.Params.FreezeFee.String())
	}
	if data.Params.UnFreezeFee.IsNegative() {
		return fmt.Errorf("invalid unfreeze fee: %s", data.Params.UnFreezeFee.String())
	}
	if data.Params.TransferOwnerFee.IsNegative() {
		return fmt.Errorf("invalid transfer owner fee: %s", data.Params.TransferOwnerFee.String())
	}
	if data.Params.DescribeFee.IsNegative() {
		return fmt.Errorf("invalid describe fee: %s", data.Params.DescribeFee.String())
	}
	return nil
}
