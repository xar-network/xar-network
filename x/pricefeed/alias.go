package pricefeed

import (
	"github.com/zar-network/zar-network/x/pricefeed/internal/keeper"
	"github.com/zar-network/zar-network/x/pricefeed/internal/types"
)

type (
	Keeper = keeper.Keeper
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
