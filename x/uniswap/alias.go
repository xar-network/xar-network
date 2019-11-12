package uniswap

import (
	"github.com/xar-network/xar-network/x/uniswap/internal/keeper"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
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
