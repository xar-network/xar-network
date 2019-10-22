package cdp

import (
	"github.com/zar-network/zar-network/x/cdp/internal/keeper"
	"github.com/zar-network/zar-network/x/cdp/internal/types"
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
