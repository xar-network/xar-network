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

package types

import (
	time1 "time"

	"github.com/tendermint/tendermint/types/time"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxTimeInForce = 600

type Order struct {
	ID                store.EntityID     `json:"id"`
	Owner             sdk.AccAddress     `json:"owner"`
	MarketID          store.EntityID     `json:"market"`
	Direction         matcheng.Direction `json:"direction"`
	Price             sdk.Uint           `json:"price"`
	Quantity          sdk.Uint           `json:"quantity"`
	TimeInForceBlocks uint16             `json:"time_in_force_blocks"`
	CreatedBlock      int64              `json:"created_block"`
	CreatedTime		  time1.Time		 `json:"created_time"`
}

func New(owner sdk.AccAddress, marketID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16, created int64) Order {
	return Order{
		Owner:             owner,
		MarketID:          marketID,
		Direction:         direction,
		Price:             price,
		Quantity:          quantity,
		TimeInForceBlocks: tif,
		CreatedBlock:      created,
		CreatedTime:	   time.Now(),
	}
}
