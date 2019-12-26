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
	"time"

	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EventHandler interface {
	OnEvent(event interface{}) error
}

type Batch struct {
	BlockNumber   int64
	BlockTime     time.Time
	MarketID      store.EntityID
	ClearingPrice sdk.Uint
	Bids          []matcheng.AggregatePrice
	Asks          []matcheng.AggregatePrice
}

type Fill struct {
	OrderID     store.EntityID
	MarketID    store.EntityID
	Owner       sdk.AccAddress
	Pair        string
	Direction   matcheng.Direction
	QtyFilled   sdk.Uint
	QtyUnfilled sdk.Uint
	BlockNumber int64
	BlockTime   int64
	Price       sdk.Uint
}

type OrderCreated struct {
	ID                store.EntityID
	Owner             sdk.AccAddress
	MarketID          store.EntityID
	Direction         matcheng.Direction
	Price             sdk.Uint
	Quantity          sdk.Uint
	TimeInForceBlocks uint16
	CreatedBlock      int64
	CreatedTime		  time.Time
}

type OrderCancelled struct {
	OrderID store.EntityID
}

type BurnCreated struct {
	ID          store.EntityID
	AssetID     store.EntityID
	BlockNumber int64
	Burner      sdk.AccAddress
	Beneficiary []byte
	Quantity    sdk.Uint
}
