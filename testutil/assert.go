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

package testutil

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AssertEqualUints(t *testing.T, a sdk.Uint, b sdk.Uint, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualInts(t *testing.T, a sdk.Int, b sdk.Int, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualEntityIDs(t *testing.T, a store.EntityID, b store.EntityID, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualHex(t *testing.T, exp string, actual []byte) {
	assert.Equal(t, exp, hex.EncodeToString(actual))
}
