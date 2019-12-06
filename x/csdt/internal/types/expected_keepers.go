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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/xar-network/xar-network/x/oracle"
)

type BankKeeper interface {
	GetCoins(sdk.Context, sdk.AccAddress) sdk.Coins
	HasCoins(sdk.Context, sdk.AccAddress, sdk.Coins) bool
	AddCoins(sdk.Context, sdk.AccAddress, sdk.Coins) (sdk.Coins, sdk.Error)
	SubtractCoins(sdk.Context, sdk.AccAddress, sdk.Coins) (sdk.Coins, sdk.Error)
}

type OracleKeeper interface {
	GetCurrentPrice(sdk.Context, string) oracle.CurrentPrice
	// These are used for testing TODO replace mockApp with keeper in tests to remove these
	AddAsset(sdk.Context, string, string, oracle.Asset) error
	SetPrice(sdk.Context, sdk.AccAddress, string, sdk.Dec, time.Time) (oracle.PostedPrice, sdk.Error)
	SetCurrentPrices(sdk.Context) sdk.Error
	SetParams(ctx sdk.Context, params oracle.Params)
}

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
	SetSupply(ctx sdk.Context, supply exported.SupplyI)
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}
