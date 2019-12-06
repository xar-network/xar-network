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

package matcheng

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/testutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNormalizeQuoteQuantity(t *testing.T) {
	tests := [][3]sdk.Uint{
		{testutil.ToBaseUnitsDecimals(10, 0), testutil.ToBaseUnitsDecimals(2, 0), testutil.ToBaseUnitsDecimals(20, 0)},
		{testutil.ToBaseUnitsDecimals(1, 0), testutil.ToBaseUnitsDecimals(10, 0), testutil.ToBaseUnitsDecimals(10, 0)},
		{testutil.ToBaseUnitsDecimals(10, 0), testutil.ToBaseUnitsDecimals(1, 3), testutil.ToBaseUnitsDecimals(10, 3)},
		{testutil.ToBaseUnitsDecimals(2, 2), testutil.ToBaseUnitsDecimals(3, 3), testutil.ToBaseUnitsDecimals(6, 5)},
		{sdk.NewUint(1), testutil.ToBaseUnitsDecimals(1, 0), sdk.NewUint(1)},
	}

	_, err := NormalizeQuoteQuantity(sdk.NewUint(1), sdk.NewUint(1))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity too small to represent")

	for _, tt := range tests {
		t.Run(fmt.Sprintf("price %s quantity %s", tt[0].String(), tt[1].String()), func(t *testing.T) {
			res, err := NormalizeQuoteQuantity(tt[0], tt[1])
			require.NoError(t, err)
			testutil.AssertEqualUints(t, tt[2], res)
		})
	}

}
