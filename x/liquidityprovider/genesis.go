package liquidityprovider

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidityprovider/internal/types"
)

type GenesisState struct {
	Accounts []types.LiquidityProviderAccount `json:"accounts" yaml:"accounts"`
}

func defaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(_ sdk.Context, _ Keeper, _ GenesisState) {

}
