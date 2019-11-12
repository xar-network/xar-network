package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSwap struct {
	Sender      sdk.AccAddress `json:"sender"`
	Asset       sdk.Coin       `json:"asset"`
	TargetDenom string         `json:"target_denom"`
}

var _ sdk.Msg = MsgSwap{}

func NewMsgSwap(sender sdk.AccAddress, asset sdk.Coin, targetDenom string) MsgSwap {
	return MsgSwap{
		Sender:      sender,
		Asset:       asset,
		TargetDenom: targetDenom,
	}
}

// Route Implements Msg
func (msg MsgSwap) Route() string { return ModuleName }

// Type Implements Msg
func (msg MsgSwap) Type() string { return "swap" }

// ValidateBasic Implements Msg.
func (msg MsgSwap) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
