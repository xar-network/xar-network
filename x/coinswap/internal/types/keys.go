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

const (
	// ModuleName is the name of the module.
	ModuleName = "coinswap"

	// RouterKey is the message route for the coinswap module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the coinswap module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the coinswap module.
	QuerierRoute = StoreKey

	MsgTypeAddLiquidity    = "add_liquidity"
	MsgTypeRemoveLiquidity = "remove_liquidity"
	MsgTypeSwapOrder       = "swap_order"
)
