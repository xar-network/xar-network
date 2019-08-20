package issue

import (
	"github.com/zar-network/zar-network/x/issue/client"
	"github.com/zar-network/zar-network/x/issue/client/cli"
	"github.com/zar-network/zar-network/x/issue/internal/keeper"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

type (
	BaseKeeper    = keeper.BaseKeeper // ibc module depends on this
	Keeper        = keeper.Keeper
	CoinIssueInfo = types.CoinIssueInfo
	Approval      = types.Approval
	IssueFreeze   = types.IssueFreeze
	Params        = types.Params
	Hooks         = keeper.Hooks
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
	ModuleCdc       = types.ModuleCdc
	NewModuleClient = client.NewModuleClient
	NewKeeper       = keeper.NewKeeper
	//GetAccountCmd   = cli.GetAccountCmd
	QueryCmd      = cli.QueryCmd
	RegisterCodec = types.RegisterCodec
	DefaultParams = types.DefaultParams
)
