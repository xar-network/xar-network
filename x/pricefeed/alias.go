/**

Baseline from Kava Cosmos Module

**/

package pricefeed

import (
	"github.com/xar-network/xar-network/x/pricefeed/internal/keeper"
	"github.com/xar-network/xar-network/x/pricefeed/internal/types"
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
