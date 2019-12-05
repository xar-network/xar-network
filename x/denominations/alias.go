package denominations

import (
	"github.com/xar-network/xar-network/x/denominations/internal/keeper"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
