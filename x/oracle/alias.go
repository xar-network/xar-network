package oracle

import (
	"github.com/xar-network/xar-network/x/oracle/internal/keeper"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

type (
	Keeper       = keeper.Keeper
	CurrentPrice = types.CurrentPrice
	PostedPrice  = types.PostedPrice
)

const (
	DefaultCodespace  = types.DefaultCodespace
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
