/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

type (
	Keeper           = keeper.Keeper
	CSDT             = types.CSDT
	CSDTs            = types.CSDTs
	Params           = types.Params
	CollateralParams = types.CollateralParams
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	StableDenom       = types.StableDenom
)

var (
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
)
