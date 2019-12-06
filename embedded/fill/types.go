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

package fill

import (
	"github.com/xar-network/xar-network/embedded"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Fill struct {
	OrderID     store.EntityID     `json:"order_id"`
	Owner       sdk.AccAddress     `json:"owner"`
	Pair        string             `json:"pair"`
	Direction   matcheng.Direction `json:"direction"`
	QtyFilled   sdk.Uint           `json:"qty_filled"`
	QtyUnfilled sdk.Uint           `json:"qty_unfilled"`
	BlockNumber int64              `json:"block_number"`
	Price       sdk.Uint           `json:"price"`
}

type QueryRequest struct {
	Owner      sdk.AccAddress
	StartBlock int64
	EndBlock   int64
}

type QueryResult struct {
	Fills []Fill
}

type RESTQueryResult struct {
	Fills []RESTFill `json:"fills"`
}

type RESTFill struct {
	BlockInclusion   embedded.BlockInclusion `json:"block_inclusion"`
	QuantityFilled   sdk.Uint                `json:"quantity_filled"`
	QuantityUnfilled sdk.Uint                `json:"quantity_unfilled"`
	Direction        matcheng.Direction      `json:"direction"`
	OrderID          store.EntityID          `json:"order_id"`
	Pair             string                  `json:"pair"`
	Price            sdk.Uint                `json:"price"`
	Owner            sdk.AccAddress          `json:"owner"`
}
