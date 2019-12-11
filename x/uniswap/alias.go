package uniswap

import (
	"github.com/xar-network/xar-network/x/uniswap/client/rest"
	"github.com/xar-network/xar-network/x/uniswap/internal/keeper"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

type (
	Keeper = keeper.Keeper
	MsgSwapOrder = types.MsgSwapOrder
	MsgAddLiquidity = types.MsgAddLiquidity
	MsgRemoveLiquidity = types.MsgRemoveLiquidity
	MsgTransactionOrder = types.MsgTransactionOrder
)

var (
	ErrInvalidDeadline         = types.ErrInvalidDeadline
	ErrNotPositive             = types.ErrNotPositive
	ErrCannotCreateReservePool = types.ErrCannotCreateReservePool
	ErrConstraintNotMet        = types.ErrConstraintNotMet
	ErrNotSupported            = types.ErrNotSupported
)

const (
	DefaultCodespace  = types.DefaultCodespace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
)

var (
	ModuleCdc      = types.ModuleCdc
	NewKeeper      = keeper.NewKeeper
	RegisterCodec  = types.RegisterCodec
	RegisterRoutes = rest.RegisterRoutes
)