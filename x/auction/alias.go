package auction

import (
	"github.com/zar-network/zar-network/x/auction/client/cli"
	"github.com/zar-network/zar-network/x/auction/internal/keeper"
	"github.com/zar-network/zar-network/x/auction/internal/types"
)

type (
	BaseKeeper = keeper.BaseKeeper // ibc module depends on this
	Keeper     = keeper.Keeper
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
)

var (
	ModuleCdc = types.ModuleCdc
	NewKeeper = keeper.NewKeeper
	//GetAccountCmd   = cli.GetAccountCmd
	QueryCmd      = cli.QueryCmd
	RegisterCodec = types.RegisterCodec
	DefaultParams = types.DefaultParams
)
