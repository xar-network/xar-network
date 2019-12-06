/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	// check that valid case work
	defaultParams := DefaultParams()
	err := ValidateParams(defaultParams)
	require.Nil(t, err)

	// all cases should return an error
	invalidTests := []struct {
		name   string
		params Params
		result bool
	}{
		{"fee == 0 ", NewParams(sdk.ZeroDec()), false},
	}

	for _, tc := range invalidTests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateParams(tc.params)
			if err != nil {
				require.False(t, tc.result)
			} else {
				require.True(t, tc.result)
			}
		})
	}
}
