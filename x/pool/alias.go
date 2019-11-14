package pool

import (
	"github.com/xar-network/xar-network/x/pool/internal/keeper"
	"github.com/xar-network/xar-network/x/pool/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamSpace = types.DefaultParamspace
	StoreKey          = types.StoreKey
)

var (
	NewKeeper = keeper.NewKeeper
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper = keeper.Keeper
)
