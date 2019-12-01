package escrow

import (
	"github.com/xar-network/xar-network/x/escrow/client/cli"
	"github.com/xar-network/xar-network/x/escrow/internal/keeper"
	"github.com/xar-network/xar-network/x/escrow/internal/types"
)

type (
	Keeper  = keeper.Keeper
	BoxInfo = types.BoxInfo
	Params  = types.Params
	Hooks   = keeper.Hooks
)

var (
	MsgCdc        = types.MsgCdc
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	SendTxCmd     = cli.SendTxCmd
	QueryCmd      = cli.QueryCmd
	WithdrawCmd   = cli.WithdrawCmd
	DefaultParams = types.DefaultParams
)

const (
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)
