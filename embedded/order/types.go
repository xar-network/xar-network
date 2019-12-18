/*

Copyright 2019 All in Bits, Inc
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

package order

import (
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	ID             store.EntityID     `json:"id"`
	Owner          sdk.AccAddress     `json:"owner"`
	MarketID       store.EntityID     `json:"market_id"`
	Direction      matcheng.Direction `json:"direction"`
	Price          sdk.Uint           `json:"price"`
	Quantity       sdk.Uint           `json:"quantity"`
	Status         string             `json:"status"`
	Type           string             `json:"type"`
	TimeInForce    uint16             `json:"time_in_force"`
	QuantityFilled sdk.Uint           `json:"quantity_filled"`
	CreatedBlock   int64              `json:"created_block"`
	CreatedTime	   int64			  `json:"created_time"`
}

type ListQueryRequest struct {
	Start          store.EntityID
	Owner          sdk.AccAddress
	Limit          int
	MarketID       []store.EntityID
	Status         []string
	UnixTimeAfter  int64
	UnixTimeBefore int64
}

type ListQueryResult struct {
	NextID store.EntityID `json:"next_id"`
	Orders []Order        `json:"orders"`
}
