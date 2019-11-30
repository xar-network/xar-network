package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgPostPrice type of PostPrice msg
	TypeMsgPostPrice = "post_price"
)

// MsgPostPrice struct representing a posted price message.
// Used by oracles to input prices to the oracle
type MsgPostPrice struct {
	From      sdk.AccAddress `json:"from" yaml:"from"`
	AssetCode string         `json:"asset_code" yaml:"asset_code"`
	Price     sdk.Dec        `json:"price" yaml:"price"`
	Expiry    sdk.Int        `json:"expiry" yaml:"expiry"`
}

// NewMsgPostPrice creates a new post price msg
func NewMsgPostPrice(
	from sdk.AccAddress,
	assetCode string,
	price sdk.Dec,
	expiry sdk.Int) MsgPostPrice {
	return MsgPostPrice{
		From:      from,
		AssetCode: assetCode,
		Price:     price,
		Expiry:    expiry,
	}
}

// Route Implements Msg.
func (msg MsgPostPrice) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgPostPrice) Type() string { return TypeMsgPostPrice }

// GetSignBytes Implements Msg.
func (msg MsgPostPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgPostPrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgPostPrice) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInternal("invalid (empty) oracle address")
	}
	if len(msg.AssetCode) == 0 {
		return sdk.ErrInternal("invalid (empty) asset code")
	}
	if msg.Price.LT(sdk.ZeroDec()) {
		return sdk.ErrInternal("invalid (negative) price")
	}
	if msg.Expiry.LT(sdk.ZeroInt()) {
		return sdk.ErrInternal("invalid (negative) expiry")
	}
	// TODO check coin denoms
	return nil
}
