package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateOrModifyCSDT creates, adds/removes collateral/stable coin from a csdt
// TODO Make this more user friendly - maybe split into four functions.
type MsgCreateOrModifyCSDT struct {
	Sender           sdk.AccAddress `json:"sender" yaml:"sender"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	CollateralChange sdk.Int        `json:"collateral_change" yaml:"collateral_change"`
	DebtChange       sdk.Int        `json:"debt_change" yaml:"debt_change"`
}

// NewMsgPlaceBid returns a new MsgPlaceBid.
func NewMsgCreateOrModifyCSDT(sender sdk.AccAddress, collateralDenom string, collateralChange sdk.Int, debtChange sdk.Int) MsgCreateOrModifyCSDT {
	return MsgCreateOrModifyCSDT{
		Sender:           sender,
		CollateralDenom:  collateralDenom,
		CollateralChange: collateralChange,
		DebtChange:       debtChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgCreateOrModifyCSDT) Route() string { return "csdt" }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgCreateOrModifyCSDT) Type() string { return "create_modify_csdt" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgCreateOrModifyCSDT) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	// TODO check coin denoms
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgCreateOrModifyCSDT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgCreateOrModifyCSDT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgTransferCSDT changes the ownership of a csdt
type MsgTransferCSDT struct {
	// TODO
}
