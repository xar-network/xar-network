package liquidityprovider

import (
	"github.com/xar-network/xar-network/x/liquidityprovider/internal/types"
)

type genesisState struct {
	Accounts []types.LiquidityProviderAccount
}

func defaultGenesisState() genesisState {
	return genesisState{}
}

//
//func InitGenesis(_ *sdk.Context,  am.keeper, gs genesisState) {
//
//}
