/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

type (
	Keeper = keeper.Keeper
	CSDT    = types.CSDT
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
