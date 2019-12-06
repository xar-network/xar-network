/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

package coinswap

import (
	"github.com/xar-network/xar-network/x/coinswap/internal/keeper"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

type (
	Keeper               = keeper.Keeper
	MsgSwapOrder         = types.MsgSwapOrder
	MsgAddLiquidity      = types.MsgAddLiquidity
	MsgRemoveLiquidity   = types.MsgRemoveLiquidity
	Params               = types.Params
	QueryLiquidityParams = types.QueryLiquidityParams
	Input                = types.Input
	Output               = types.Output
)

var (
	DefaultParamSpace = types.DefaultParamSpace
	QueryLiquidity    = types.QueryLiquidity

	RegisterCodec = types.RegisterCodec

	NewMsgSwapOrder       = types.NewMsgSwapOrder
	NewMsgAddLiquidity    = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity = types.NewMsgRemoveLiquidity
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier

	ErrInvalidDeadline  = types.ErrInvalidDeadline
	ErrNotPositive      = types.ErrNotPositive
	ErrConstraintNotMet = types.ErrConstraintNotMet

	GetUniId                    = types.GetUniId
	GetCoinMinDenomFromUniDenom = types.GetCoinMinDenomFromUniDenom
	GetUniDenom                 = types.GetUniDenom
	GetUniCoinType              = types.GetUniCoinType
	CheckUniDenom               = types.CheckUniDenom
	CheckUniId                  = types.CheckUniId
)

const (
	DefaultCodespace   = types.DefaultCodespace
	ModuleName         = types.ModuleName
	FormatUniABSPrefix = types.FormatUniABSPrefix
)
