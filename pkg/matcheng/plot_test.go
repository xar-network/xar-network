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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xar-network/xar-network/testutil/testflags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPlotCurves(t *testing.T) {
	testflags.UnitTest(t)
	expected := `"Ask"
2 0
2 10
3 10
3 20
4 20
4 30


"Bid"
3 0
3 10
2 10
2 20
1 20
1 30
0 30
`

	res := &MatchResults{
		BidAggregates: []AggregatePrice{
			{sdk.NewUint(1), sdk.NewUint(30)},
			{sdk.NewUint(2), sdk.NewUint(20)},
			{sdk.NewUint(3), sdk.NewUint(10)},
		},
		AskAggregates: []AggregatePrice{
			{sdk.NewUint(2), sdk.NewUint(10)},
			{sdk.NewUint(3), sdk.NewUint(20)},
			{sdk.NewUint(4), sdk.NewUint(30)},
		},
	}

	actual := PlotCurves(res.BidAggregates, res.AskAggregates)
	assert.Equal(t, expected, actual)

}
