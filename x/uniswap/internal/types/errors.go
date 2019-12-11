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

package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolAlreadyExists    sdk.CodeType = 101
	CodeEqualDenom                  sdk.CodeType = 102
	CodeInvalidDeadline             sdk.CodeType = 103
	CodeNotPositive                 sdk.CodeType = 104
	CodeConstraintNotMet            sdk.CodeType = 105
	CodeNotSupported                sdk.CodeType = 106
	CodeCannotCreateReservePool     sdk.CodeType = 107
	CodeInvalidAccountAddr          sdk.CodeType = 108
	CodeInvalidAccountPemission     sdk.CodeType = 109
	CodeQueryParamIsInvalid         sdk.CodeType = 110
	CodeInsufficientLiquidityAmount sdk.CodeType = 111
	CodeReservePoolNotFound         sdk.CodeType = 111
)

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
)

func MsgAccPermissionsError(moduleName string) string {
	return fmt.Sprintf("module account %s does not have permissions to burn tokens", moduleName)
}

func MsgReservePoolNotFound(moduleName string) string {
	return fmt.Sprintf("error retrieving reserve pool for ModuleAccoint name: %s", moduleName)
}

func ErrReservePoolNotFound(codespace sdk.CodespaceType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeReservePoolNotFound, MsgReservePoolNotFound(moduleName))
}

func ErrReservePoolAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeReservePoolAlreadyExists, msg)
	}
	return sdk.NewError(codespace, CodeReservePoolAlreadyExists, ReservePoolAlreadyExists)
}

func ErrEqualDenom(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeEqualDenom, msg)
	}
	return sdk.NewError(codespace, CodeEqualDenom, EqualDenom)
}

func ErrInvalidAccountAddr(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAccountAddr, msg)
}

func ErrInvalidAccountPermission(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAccountPemission, msg)
}

func ErrQueryParamIsInvalid(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeQueryParamIsInvalid, msg)
}

func ErrInvalidDeadline(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeInvalidDeadline, msg)
	}
	return sdk.NewError(codespace, CodeInvalidDeadline, InvalidDeadline)
}

func ErrNotPositive(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeNotPositive, msg)
	}
	return sdk.NewError(codespace, CodeNotPositive, AmountIsNotPositive)
}

func ErrConstraintNotMet(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeConstraintNotMet, msg)
	}
	return sdk.NewError(codespace, CodeConstraintNotMet, ConstraintNotMet)
}

func ErrNotSupported(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeNotSupported, msg)
	}
	return sdk.NewError(codespace, CodeNotSupported, NotCurrentlySupported)
}

func ErrCannotCreateReservePool(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeCannotCreateReservePool, CannotCreateReservePool)
}

func ErrInsufficientLiquidityAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientLiquidityAmount, InsufficientLiquidityAmount)
}
