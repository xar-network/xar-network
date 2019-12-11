/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package uniswap

import (
	"github.com/xar-network/xar-network/x/uniswap/client/rest"
	"github.com/xar-network/xar-network/x/uniswap/internal/keeper"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

type (
	Keeper              = keeper.Keeper
	MsgSwapOrder        = types.MsgSwapOrder
	MsgAddLiquidity     = types.MsgAddLiquidity
	MsgRemoveLiquidity  = types.MsgRemoveLiquidity
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
