/*

Copyright 2016 All in Bits, Inc
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

// test ValidateBasic for MsgSwapOrder
func TestMsgSwapOrder(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgSwapOrder
		expectPass bool
	}{
		{"no input coin", NewMsgSwapOrder(sdk.Coin{}, output, sender, recipient, true), false},
		{"zero input coin", NewMsgSwapOrder(sdk.NewCoin(denom0, sdk.ZeroInt()), output, sender, recipient, true), false},
		{"no output coin", NewMsgSwapOrder(input, sdk.Coin{}, sender, recipient, false), false},
		{"zero output coin", NewMsgSwapOrder(input, sdk.NewCoin(denom1, sdk.ZeroInt()), sender, recipient, true), false},
		{"swap and coin denomination are equal", NewMsgSwapOrder(input, sdk.NewCoin(denom0, amt), sender, recipient, true), false},
		{"no sender", NewMsgSwapOrder(input, output, emptyAddr, recipient, true), false},
		{"no recipient", NewMsgSwapOrder(input, output, sender, emptyAddr, true), false},
		{"valid MsgSwapOrder", NewMsgSwapOrder(input, output, sender, recipient, true), true},
		{"sender and recipient are same", NewMsgSwapOrder(input, output, sender, sender, true), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// test ValidateBasic for MsgAddLiquidity
func TestMsgAddLiquidity(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgAddLiquidity
		expectPass bool
	}{
		{"no deposit coin", NewMsgAddLiquidity(sdk.Coin{}, amt, sdk.OneInt(), sender), false},
		{"zero deposit coin", NewMsgAddLiquidity(sdk.NewCoin(denom1, sdk.ZeroInt()), amt, sdk.OneInt(), sender), false},
		{"invalid withdraw amount", NewMsgAddLiquidity(input, sdk.ZeroInt(), sdk.OneInt(), sender), false},
		{"invalid minumum reward bound", NewMsgAddLiquidity(input, amt, sdk.ZeroInt(), sender), false},
		{"empty sender", NewMsgAddLiquidity(input, amt, sdk.OneInt(), emptyAddr), false},
		{"valid MsgAddLiquidity", NewMsgAddLiquidity(input, amt, sdk.OneInt(), sender), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// test ValidateBasic for MsgRemoveLiquidity
func TestMsgRemoveLiquidity(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgRemoveLiquidity
		expectPass bool
	}{
		{"no withdraw coin", NewMsgRemoveLiquidity(sdk.Coin{}, amt, sdk.OneInt(), sender), false},
		{"zero withdraw coin", NewMsgRemoveLiquidity(sdk.NewCoin(denom1, sdk.ZeroInt()), amt, sdk.OneInt(), sender), false},
		{"invalid deposit amount", NewMsgRemoveLiquidity(input, sdk.ZeroInt(), sdk.OneInt(), sender), false},
		{"invalid minimum native bound", NewMsgRemoveLiquidity(input, amt, sdk.ZeroInt(), sender), false},
		{"empty sender", NewMsgRemoveLiquidity(input, amt, sdk.OneInt(), emptyAddr), false},
		{"valid MsgRemoveLiquidity", NewMsgRemoveLiquidity(input, amt, sdk.OneInt(), sender), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}

}
