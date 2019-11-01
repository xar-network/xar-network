package liquidityprovider

import (
	"github.com/xar-network/xar-network/x/liquidityprovider/keeper"
	"github.com/xar-network/xar-network/x/liquidityprovider/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
)

type (
	Keeper  = keeper.Keeper
	Account = types.LiquidityProviderAccount
)
