/**

Baseline from Kava Cosmos Module

**/

package pricefeed

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/pricefeed/internal/types"
)

// GenesisState state at gensis
type GenesisState struct {
	Assets  []types.Asset
	Oracles []types.Oracle
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	for _, asset := range genState.Assets {
		keeper.AddAsset(ctx, asset.AssetCode, asset.Description)
	}

	for _, oracle := range genState.Oracles {
		keeper.AddOracle(ctx, oracle.OracleAddress)
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		[]types.Asset{
			{AssetCode: "btc", Description: "Bitcoin"},
			{AssetCode: "bnb", Description: "Binance Chain Coin"},
			{AssetCode: "eth", Description: "Ethereum"},
			{AssetCode: "ftm", Description: "Fantom"}},
		[]types.Oracle{}}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// TODO
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	// TODO implement this
	return DefaultGenesisState()
}
