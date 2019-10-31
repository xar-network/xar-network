/**

Baseline from Kava Cosmos Module

**/

package cdp

import (
	"github.com/xar-network/xar-network/x/cdp/internal/keeper"
	"github.com/xar-network/xar-network/x/cdp/internal/types"
)

type (
	Keeper = keeper.Keeper
	CDP    = types.CDP
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
