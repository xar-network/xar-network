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

package store

import (
	"fmt"
	"strings"

	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/store/types"
)

const TablePrefix = "t"

type Table struct {
	db     dbm.DB
	prefix string
}

func NewTable(db dbm.DB, prefix string) *Table {
	return &Table{
		db:     db,
		prefix: fmt.Sprintf("%s/%s", TablePrefix, prefix),
	}
}

func (t *Table) Get(key []byte) []byte {
	return t.db.Get(PrefixKeyString(t.prefix, key))
}

func (t *Table) Has(key []byte) bool {
	return t.db.Has(PrefixKeyString(t.prefix, key))
}

func (t *Table) Set(key, value []byte) {
	t.db.Set(PrefixKeyString(t.prefix, key), value)
}

func (t *Table) Delete(key []byte) {
	t.db.Delete(PrefixKeyString(t.prefix, key))
}

func (t *Table) Iterator(start []byte, end []byte, cb IteratorCB) {
	iter := t.db.Iterator(PrefixKeyString(t.prefix, start), PrefixKeyString(t.prefix, end))
	t.iterate(iter, cb)
}

func (t *Table) ReverseIterator(start []byte, end []byte, cb IteratorCB) {
	iter := t.db.ReverseIterator(PrefixKeyString(t.prefix, start), PrefixKeyString(t.prefix, end))
	t.iterate(iter, cb)
}

func (t *Table) PrefixIterator(start []byte, cb IteratorCB) {
	start = PrefixKeyString(t.prefix, start)
	iter := t.db.Iterator(start, sdk.PrefixEndBytes(start))
	t.iterate(iter, cb)
}

func (t *Table) ReversePrefixIterator(start []byte, cb IteratorCB) {
	start = PrefixKeyString(t.prefix, start)
	iter := t.db.ReverseIterator(start, sdk.PrefixEndBytes(start))
	t.iterate(iter, cb)
}

func (t *Table) iterate(iter dbm.Iterator, cb IteratorCB) {
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		k := []byte(strings.TrimPrefix(string(iter.Key()), t.prefix+"/"))
		v := iter.Value()

		if !cb(k, v) {
			return
		}
	}
}

func (t *Table) Substore(prefix string) ArchiveStore {
	return &Table{
		db:     t.db,
		prefix: fmt.Sprintf("%s/%s", t.prefix, prefix),
	}
}
