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

package serde

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/testutil/testflags"
)

type example struct {
	Value HexBytes `json:"value"`
}

func TestHexBytes(t *testing.T) {
	testflags.UnitTest(t)
	ser := example{Value: []byte{0x99}}
	out, err := json.Marshal(ser)
	require.NoError(t, err)
	assert.Equal(t, "{\"value\":\"0x99\"}", string(out))
	var deser example
	err = json.Unmarshal(out, &deser)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(ser.Value, deser.Value))

	ser = example{Value: nil}
	out, err = json.Marshal(ser)
	require.NoError(t, err)
	assert.Equal(t, "{\"value\":null}", string(out))
	err = json.Unmarshal(out, &deser)
	require.NoError(t, err)
	assert.Nil(t, deser.Value)

	err = json.Unmarshal([]byte("{\"value\":\"}"), &deser)
	assert.Error(t, err)
	err = json.Unmarshal([]byte("{\"value\":\"\"}"), &deser)
	assert.Error(t, err)
}
