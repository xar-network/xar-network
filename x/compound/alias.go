package compound

import (
	"github.com/xar-network/xar-network/x/compound/internal/keeper"
	"github.com/xar-network/xar-network/x/compound/internal/types"
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
