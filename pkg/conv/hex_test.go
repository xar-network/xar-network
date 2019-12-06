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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/testutil/testflags"
)

func TestHexToBytes(t *testing.T) {
	testflags.UnitTest(t)
	a, err := HexToBytes("0x0101")
	require.NoError(t, err)
	b, err := HexToBytes("0101")
	require.NoError(t, err)
	assert.EqualValues(t, a, b)
	assert.Equal(t, []byte{0x01, 0x01}, a)

	_, err = HexToBytes("foo")
	assert.Error(t, err)
}
