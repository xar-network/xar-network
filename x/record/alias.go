package record

import (
	"github.com/xar-network/xar-network/x/record/client"
	"github.com/xar-network/xar-network/x/record/internal/keeper"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

type (
	Keeper     = keeper.Keeper
	RecordInfo = types.RecordInfo
)

var (
	NewKeeper       = keeper.NewKeeper
	NewModuleClient = client.NewModuleClient
	RegisterCodec   = types.RegisterCodec
	ModuleCdc       = types.ModuleCdc
)

const (
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)
