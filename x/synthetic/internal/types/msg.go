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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBuySynthetic purchases a synthetic position
type MsgBuySynthetic struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Coin   sdk.Coin       `json:"coin" yaml:"coin"`
}

// NewMsgBuySynthetic returns a new MsgBuySynthetic.
func NewMsgBuySynthetic(sender sdk.AccAddress, coin sdk.Coin) MsgBuySynthetic {
	return MsgBuySynthetic{
		Sender: sender,
		Coin:   coin,
	}
}

// Route return the message type used for routing the message.
func (msg MsgBuySynthetic) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgBuySynthetic) Type() string { return "buy_synthetic" }

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgBuySynthetic) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.Coin.IsZero() || !msg.Coin.IsValid() || msg.Coin.IsNegative() {
		return sdk.ErrInternal("invalid (empty) coin")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgBuySynthetic) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgBuySynthetic) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgSellSynthetic purchases a synthetic position
type MsgSellSynthetic struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Coin   sdk.Coin       `json:"coin" yaml:"coin"`
}

// NewMsgSellSynthetic returns a new MsgSellSynthetic.
func NewMsgSellSynthetic(sender sdk.AccAddress, coin sdk.Coin) MsgSellSynthetic {
	return MsgSellSynthetic{
		Sender: sender,
		Coin:   coin,
	}
}

// Route return the message type used for routing the message.
func (msg MsgSellSynthetic) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgSellSynthetic) Type() string { return "sell_synthetic" }

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgSellSynthetic) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.Coin.IsZero() || !msg.Coin.IsValid() || msg.Coin.IsNegative() {
		return sdk.ErrInternal("invalid (empty) coin")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgSellSynthetic) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgSellSynthetic) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
