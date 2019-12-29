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
	sdk "github.com/cosmos/cosmos-sdk/types"
	supply "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) supply.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc supply.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error

	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}
