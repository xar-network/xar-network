/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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
	// DefaultCodespace codespace for the module
	DefaultCodespace sdk.CodespaceType = ModuleName

	// CodeEmptyInput error code for empty input errors
	CodeEmptyInput sdk.CodeType = 1
	// CodeExpired error code for expired prices
	CodeExpired sdk.CodeType = 2
	// CodeInvalidPrice error code for all input prices expired
	CodeInvalidPrice sdk.CodeType = 3
	// CodeInvalidAsset error code for invalid asset
	CodeInvalidAsset sdk.CodeType = 4
	// CodeInvalidOracle error code for invalid oracle
	CodeInvalidOracle sdk.CodeType = 5
)

// ErrEmptyInput Error constructor
func ErrEmptyInput(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyInput, fmt.Sprintf("Input must not be empty."))
}

// ErrExpired Error constructor for posted price messages with expired price
func ErrExpired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeExpired, fmt.Sprintf("Price is expired."))
}

// ErrNoValidPrice Error constructor for posted price messages with expired price
func ErrNoValidPrice(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrice, fmt.Sprintf("All input prices are expired."))
}

// ErrInvalidAsset Error constructor for posted price messages for invalid assets
func ErrInvalidAsset(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAsset, fmt.Sprintf("Asset code does not exist."))
}

// ErrExistingAsset Error constructor for posted price messages for invalid assets
func ErrExistingAsset(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAsset, fmt.Sprintf("Asset code exists."))
}

// ErrInvalidOracle Error constructor for posted price messages for invalid oracles
func ErrInvalidOracle(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOracle, fmt.Sprintf("Oracle does not exist or not authorized."))
}
