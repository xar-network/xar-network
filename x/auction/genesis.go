/**

Baseline from Kava Cosmos Module

**/

package auction

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

// GenesisAuctions type for an array of auctions
type GenesisAuctions []types.Auction

// GenesisState - auction state that must be provided at genesis
type GenesisState struct {
	AuctionParams types.AuctionParams `json:"auction_params" yaml:"auction_params"`
	Auctions      GenesisAuctions     `json:"genesis_auctions" yaml:"genesis_auctions"`
}

// NewGenesisState returns a new genesis state object for auctions module
func NewGenesisState(ap types.AuctionParams, ga GenesisAuctions) GenesisState {
	return GenesisState{
		AuctionParams: ap,
		Auctions:      ga,
	}
}

// DefaultGenesisState defines default genesis state for auction module
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultAuctionParams(), GenesisAuctions{})
}

// Equal checks whether two GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis validates genesis inputs. Returns error if validation of any input fails.
func ValidateGenesis(data GenesisState) error {
	if err := data.AuctionParams.Validate(); err != nil {
		return err
	}
	return nil
}

// InitGenesis - initializes the store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.AuctionParams)

	for _, a := range data.Auctions {
		keeper.SetAuction(ctx, a)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	var genAuctions GenesisAuctions
	iterator := keeper.GetAuctionIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {

		auction := keeper.DecodeAuction(ctx, iterator.Value())
		genAuctions = append(genAuctions, auction)

	}
	return NewGenesisState(params, genAuctions)
}
