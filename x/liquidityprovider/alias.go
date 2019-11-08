package liquidityprovider

import (
	"github.com/xar-network/xar-network/x/liquidityprovider/internal/keeper"
	"github.com/xar-network/xar-network/x/liquidityprovider/internal/types"
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
