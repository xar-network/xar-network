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

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFormatCoin(t *testing.T) {
	out := FormatCoin(NewEntityID(1), sdk.NewUint(100000))
	assert.True(t, out.Amount.Equal(sdk.NewInt(100000)))
	assert.Equal(t, "asset1", out.Denom)
}

func TestFormatDenom(t *testing.T) {
	assert.Equal(t, "asset99", FormatDenom(NewEntityID(99)))
}
