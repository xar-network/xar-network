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

package price

import (
	"time"

	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/store/types"
)

type IteratorCB func(tick Tick) bool

type Keeper struct {
	as  store.ArchiveStore
	cdc *codec.Codec
}

func NewKeeper(db dbm.DB, cdc *codec.Codec) Keeper {
	return Keeper{
		as:  store.NewTable(db, EntityName),
		cdc: cdc,
	}
}

func (k Keeper) ReverseIteratorByMarket(mktID store.EntityID, cb IteratorCB) {
	k.as.PrefixIterator(tickIterKey(mktID), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) ReverseIteratorByMarketFrom(mktID store.EntityID, from time.Time, cb IteratorCB) {
	k.as.ReverseIterator(tickKey(mktID, 0), sdk.PrefixEndBytes(tickKey(mktID, 0)), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) IteratorByMarketAndInterval(mktID store.EntityID, from time.Time, to time.Time, cb IteratorCB) {
	k.as.Iterator(tickKey(mktID, from.Unix()), sdk.PrefixEndBytes(tickKey(mktID, to.Unix())), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) OnFillEvent(event types.Fill) {
	tick := Tick{
		MarketID:    event.MarketID,
		Pair:        event.Pair,
		BlockNumber: event.BlockNumber,
		BlockTime:   event.BlockTime,
		Price:       event.Price,
	}
	storedB := k.cdc.MustMarshalBinaryBare(tick)
	k.as.Set(tickKey(event.MarketID, tick.BlockTime), storedB)
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Fill:
		k.OnFillEvent(ev)
	}

	return nil
}

func tickKey(mktID store.EntityID, blockTime int64) []byte {
	return store.PrefixKeyBytes(tickIterKey(mktID), store.Int64Subkey(blockTime))
}

func tickIterKey(mktID store.EntityID) []byte {
	return store.PrefixKeyString("tick", mktID.Bytes())
}
