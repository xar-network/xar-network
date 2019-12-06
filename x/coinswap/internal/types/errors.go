// Package types nolint
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

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolNotExists         sdk.CodeType = 101
	CodeEqualDenom                   sdk.CodeType = 102
	CodeInvalidDeadline              sdk.CodeType = 103
	CodeNotPositive                  sdk.CodeType = 104
	CodeConstraintNotMet             sdk.CodeType = 105
	CodeIllegalDenom                 sdk.CodeType = 106
	CodeIllegalUniId                 sdk.CodeType = 107
	CodeReservePoolInsufficientFunds sdk.CodeType = 108
)

func ErrReservePoolNotExists(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, "reserve pool not exists")
}

func ErrEqualDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeEqualDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeEqualDenom, "input and output denomination are equal")
}

func ErrIllegalUniId(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeIllegalUniId, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeIllegalUniId, "illegal liquidity id")
}

func ErrIllegalDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeIllegalDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeIllegalDenom, "illegal denomination")
}

func ErrInvalidDeadline(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, "invalid deadline")
}

func ErrNotPositive(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeNotPositive, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeNotPositive, "amount is not positive")
}

func ErrConstraintNotMet(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, "constraint not met")
}

func ErrInsufficientFunds(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, "constraint not met")
}
