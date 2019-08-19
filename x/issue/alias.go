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

var (
	ModuleCdc       = types.ModuleCdc
	NewModuleClient = client.NewModuleClient
	//GetAccountCmd   = cli.GetAccountCmd
	QueryCmd      = cli.QueryCmd
	RegisterCodec = types.RegisterCodec
	DefaultParams = types.DefaultParams
)

const (
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)
