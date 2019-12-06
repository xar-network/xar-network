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

package conv

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xar-network/xar-network/testutil/testflags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestUint2Big(t *testing.T) {
	testflags.UnitTest(t)
	a := sdk.NewUint(1)
	b := big.NewInt(1)
	assert.Equal(t, "1", SDKUint2Big(a).String())
	assert.EqualValues(t, b, SDKUint2Big(a))
}
