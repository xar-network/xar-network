package uniswap

import (
	"github.com/xar-network/xar-network/x/uniswap/internal/keeper"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

type (
	Keeper             = keeper.Keeper
	MsgSwapOrder       = types.MsgSwapOrder
	MsgAddLiquidity    = types.MsgAddLiquidity
	MsgRemoveLiquidity = types.MsgRemoveLiquidity
)

var (
	ErrInvalidDeadline  = types.ErrInvalidDeadline
	ErrNotPositive      = types.ErrNotPositive
	ErrConstraintNotMet = types.ErrConstraintNotMet
	ErrNotSupported     = types.ErrNotSupported
)

const (
	DefaultCodespace = types.DefaultCodespace
	ModuleName       = types.ModuleName
)
