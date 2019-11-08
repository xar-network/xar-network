package authority

import (
	"github.com/xar-network/xar-network/x/authority/internal/keeper"
	"github.com/xar-network/xar-network/x/authority/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

type (
	Keeper = keeper.Keeper
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
)
