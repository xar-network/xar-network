// nolint
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

package pool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeIncorrectNativeDenom sdk.CodeType = 101 + iota
	CodeIncorrectNonNativeDenom
	CodeIncorrectNativeAmount
	CodeIncorrectNonNativeAmount
	CodeNoNativeDenomPresent
	CodeNotAllDenomsAreInPool
)

var codespace sdk.CodespaceType = "uniswap"

func SetCodespace(c sdk.CodespaceType) {
	codespace = c
}

// constant set for error messages
const (
	ReservePoolAlreadyExists         = "reserve pool already exists"
	EqualDenom                       = "input and output denomination are equal"
	InvalidDeadline                  = "invalid deadline"
	AmountIsNotPositive              = "amount is not positive"
	ConstraintNotMet                 = "constraint not met"
	NotCurrentlySupported            = "not currently supported"
	InvalidQueryParameter            = "query parameter is invalid"
	CannotCreateReservePool          = "cannot create reserve pool"
	InsufficientLiquidityAmount      = "insufficient liquidity amount"
	InsufficientCoins                = "sender does not have sufficient funds"
	LiquidityAddDeadLineHasPassed    = "deadline has passed for MsgAddLiquidity"
	LiquidityRemoveDeadLineHasPassed = "deadline has passed for MsgRemoveLiquidity"
	IncorrectNativeDenomMsg          = "native coin denom from add liquidity request does not equal to native coin denom from a reserve pool"
	IncorrectNonNativeDenomMsg       = "non-native coin denom from add liquidity request does not equal to non-native coin denom from a reserve pool"
	NoNativeDenomPresentMsg          = "native denom is not present"
	NotAllDenomsAreInPoolMsg         = "native denom is not present"
)

var ErrIncorrectNativeDenom = sdk.NewError(codespace, CodeIncorrectNativeDenom, IncorrectNativeDenomMsg)
var ErrIncorrectNonNativeDenom = sdk.NewError(codespace, CodeIncorrectNonNativeDenom, IncorrectNonNativeDenomMsg)
var ErrNoNativeDenomPresent = sdk.NewError(codespace, CodeNoNativeDenomPresent, NoNativeDenomPresentMsg)
var ErrNotAllDenomsAreInPool = sdk.NewError(codespace, CodeNotAllDenomsAreInPool, NotAllDenomsAreInPoolMsg)

func ErrIncorrectNativeAmount(msg string) sdk.Error {
	return sdk.NewError(codespace, CodeIncorrectNativeAmount, msg)
}

func ErrIncorrectNonNativeAmount(msg string) sdk.Error {
	return sdk.NewError(codespace, CodeIncorrectNonNativeAmount, msg)
}
