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
	"testing"

	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/testutil/testflags"

	"github.com/cosmos/cosmos-sdk/codec"
)

type incrementingTest struct {
	ID  EntityID
	Foo string
	Bar int
}

func (it *incrementingTest) GetID() EntityID {
	return it.ID
}

func (it *incrementingTest) SetID(id EntityID) {
	it.ID = id
}

func TestIncrementing(t *testing.T) {
	testflags.UnitTest(t)
	db := dbm.NewMemDB()
	inc := NewIncrementing(db, codec.New())
	data := incrementingTest{
		ID:  NewEntityID(1),
		Foo: "hello",
		Bar: 1,
	}

	err := inc.Insert(&data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be zero")

	data.ID = NewEntityID(0)
	err = inc.Insert(&data)
	assert.NoError(t, err)

	var retrieved incrementingTest
	err = inc.ByID(data.ID, &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "hello", retrieved.Foo)
	assert.Equal(t, 1, retrieved.Bar)
	assert.True(t, NewEntityID(1).Equals(retrieved.ID))
	assert.True(t, inc.HasID(retrieved.ID))

	err = inc.ByID(NewEntityID(999), &retrieved)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	data.ID = NewEntityID(0)
	err = inc.Insert(&data)
	assert.NoError(t, err)
	expID := NewEntityID(2)
	assert.True(t, inc.HasID(expID))
	assert.True(t, expID.Equals(inc.HeadID()))
}
