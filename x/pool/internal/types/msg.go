package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDepositFund inserts a given amount into the pool
type MsgDepositFund struct {
	Sender sdk.AccAddress
	Amount sdk.Coin
}

func NewMsgDepositFund(sender sdk.AccAddress, amount sdk.Coin) MsgDepositFund {
	return MsgDepositFund{
		Sender: sender,
		Amount: amount,
	}
}

// Route return the message type used for routing the message.
func (msg MsgDepositFund) Route() string { return "pool" }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgDepositFund) Type() string { return "deposit_fund" }

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgDepositFund) ValidateBasic() sdk.Error {

	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgDepositFund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgDepositFund) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ============================================

// MsgWithdrawFund inserts a given amount into the pool
type MsgWithdrawFund struct {
	Sender sdk.AccAddress
	Amount sdk.Coin
}

func NewMsgWithdrawFund(sender sdk.AccAddress, amount sdk.Coin) MsgWithdrawFund {
	return MsgWithdrawFund{
		Sender: sender,
		Amount: amount,
	}
}

// Route return the message type used for routing the message.
func (msg MsgWithdrawFund) Route() string { return "pool" }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgWithdrawFund) Type() string { return "withdraw_fund" }

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgWithdrawFund) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	// TODO check coin denoms
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgWithdrawFund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgWithdrawFund) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
