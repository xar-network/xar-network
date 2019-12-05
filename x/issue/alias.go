/**

Baseline from Hashgard Cosmos Module

**/

package issue

import (
	"github.com/xar-network/xar-network/x/issue/client/cli"
	"github.com/xar-network/xar-network/x/issue/internal/keeper"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

type (
	Keeper        = keeper.Keeper
	CoinIssueInfo = types.CoinIssueInfo
	Approval      = types.Approval
	IssueFreeze   = types.IssueFreeze
	Params        = types.Params
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
	ModuleCdc = types.ModuleCdc
	NewKeeper = keeper.NewKeeper
	//GetAccountCmd   = cli.GetAccountCmd
	QueryCmd      = cli.QueryCmd
	RegisterCodec = types.RegisterCodec
	DefaultParams = types.DefaultParams
)
