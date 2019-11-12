package uniswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	CoinDenom string       `json:"coin_denom"`
	Pools     []types.Pool `json:"pools`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(coinDenom string, pools []types.Pool) GenesisState {
	return GenesisState{
		CoinDenom: coinDenom,
		Pools:     pools,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		"ftm", []types.Pool{
			types.Pool{
				BalanceCoin:  sdk.NewCoin("ftm", sdk.NewInt(76)),
				BalanceToken: sdk.NewCoin("csdt", sdk.NewInt(1)),
			},
		},
	)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	config := types.PoolConfig{
		CoinDenom: data.CoinDenom,
	}

	keeper.SetPoolConfig(ctx, config)

	for _, pool := range data.Pools {
		err := keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	panic("not yet implemented")
	return GenesisState{}
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }
