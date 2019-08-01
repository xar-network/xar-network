package issue

import (
	"github.com/zar-network/zar-network/x/issue/client"
	"github.com/zar-network/zar-network/x/issue/client/cli"
	"github.com/zar-network/zar-network/x/issue/config"
	"github.com/zar-network/zar-network/x/issue/internal/keeper"
	"github.com/zar-network/zar-network/x/issue/internal/types"
	"github.com/zar-network/zar-network/x/issue/msgs"
)

type (
	Keeper        = keeper.Keeper
	CoinIssueInfo = types.CoinIssueInfo
	Approval      = types.Approval
	IssueFreeze   = types.IssueFreeze
	Params        = config.Params
	Hooks         = keeper.Hooks
)

var (
	MsgCdc          = msgs.MsgCdc
	NewKeeper       = keeper.NewKeeper
	NewModuleClient = client.NewModuleClient
	//GetAccountCmd   = cli.GetAccountCmd
	QueryCmd      = cli.QueryCmd
	RegisterCodec = msgs.RegisterCodec
	DefaultParams = config.DefaultParams
)

const (
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)
