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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xar-network/xar-network/testutil/testflags"

	"github.com/cosmos/cosmos-sdk/types"
)

type lastCall struct {
	key   []byte
	start []byte
	end   []byte
}

type dumbKVStore struct {
	last lastCall
}

func (d *dumbKVStore) Get(key []byte) []byte {
	d.last = lastCall{key: key}
	return nil
}

func (d *dumbKVStore) Has(key []byte) bool {
	d.last = lastCall{key: key}
	return true
}

func (d *dumbKVStore) Set(key, value []byte) {
	d.last = lastCall{key: key}
}

func (d *dumbKVStore) Delete(key []byte) {
	d.last = lastCall{key: key}
}

func (d *dumbKVStore) Iterator(start, end []byte) types.Iterator {
	d.last = lastCall{start: start, end: end}
	return nil
}

func (d *dumbKVStore) ReverseIterator(start, end []byte) types.Iterator {
	d.last = lastCall{start: start, end: end}
	return nil
}

func TestPrefixed(t *testing.T) {
	testflags.UnitTest(t)
	kvs := &dumbKVStore{}
	pref := NewPrefixed(kvs, []byte{0x44})

	pref.Get([]byte{0x01})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x01}, kvs.last.key))

	pref.Has([]byte{0x02})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x02}, kvs.last.key))

	pref.Set([]byte{0x03}, nil)
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x03}, kvs.last.key))

	pref.Delete([]byte{0x04})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x04}, kvs.last.key))

	pref.Iterator([]byte{0x05}, []byte{0x06})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x05}, kvs.last.start))
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x06}, kvs.last.end))

	pref.ReverseIterator([]byte{0x07}, []byte{0x08})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x07}, kvs.last.start))
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x08}, kvs.last.end))
}
