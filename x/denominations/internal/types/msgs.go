package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgMint{}
var _ sdk.Msg = MsgBurn{}

type MsgMint struct {
	Coins sdk.Coins
	From  sdk.AccAddress
}

type MsgBurn struct {
	Coins sdk.Coins
	From  sdk.AccAddress
}

func NewMsgMint(from sdk.AccAddress, coins sdk.Coins) MsgMint {
	return MsgMint{
		Coins: coins,
		From:  from,
	}
}

func NewMsgBurn(from sdk.AccAddress, coins sdk.Coins) MsgBurn {
	return MsgBurn{
		Coins: coins,
		From:  from,
	}
}
func (msg MsgBurn) Route() string { return RouterKey }

func (msg MsgBurn) Type() string { return "burn" }

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}

	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins(msg.Coins.String())
	}

	return nil
}

func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgMint) Route() string { return RouterKey }

func (msg MsgMint) Type() string { return "mint" }

func (msg MsgMint) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}

	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins(msg.Coins.String())
	}

	return nil
}

func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
