package auction

import (
	"github.com/zar-network/zar-network/x/auction/internal/keeper"
	"github.com/zar-network/zar-network/x/auction/internal/types"
)

type (
	Keeper = keeper.Keeper
	ID     = types.ID
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
