package uniswap

import (
	"github.com/zar-network/zar-network/x/uniswap/internal/keeper"
	"github.com/zar-network/zar-network/x/uniswap/internal/types"
)

type (
	Keeper       = keeper.Keeper
	CurrentPrice = types.CurrentPrice
	PostedPrice  = types.PostedPrice
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
)

var (
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
)
