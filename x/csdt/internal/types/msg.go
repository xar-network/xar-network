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

// NewMsgCreateOrModifyCSDT returns a new MsgCreateOrModifyCSDT.
func NewMsgCreateOrModifyCSDT(sender sdk.AccAddress, collateralDenom string, collateralChange sdk.Int, debtChange sdk.Int) MsgCreateOrModifyCSDT {
	return MsgCreateOrModifyCSDT{
		Sender:           sender,
		CollateralDenom:  collateralDenom,
		CollateralChange: collateralChange,
		DebtChange:       debtChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgCreateOrModifyCSDT) Route() string { return ModuleName }

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

// MsgDepositCollateral adds collateral to CSDT.
type MsgDepositCollateral struct {
	Sender           sdk.AccAddress `json:"sender" yaml:"sender"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	CollateralChange sdk.Int        `json:"collateral_change" yaml:"collateral_change"`
}

// NewMsgDepositCollateral returns a new MsgDepositCollateral.
func NewMsgDepositCollateral(sender sdk.AccAddress, collateralDenom string, collateralChange sdk.Int) MsgDepositCollateral {
	return MsgDepositCollateral{
		Sender:           sender,
		CollateralDenom:  collateralDenom,
		CollateralChange: collateralChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgDepositCollateral) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgDepositCollateral) Type() string { return "deposit_collateral" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgDepositCollateral) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.CollateralChange.IsNegative() || msg.CollateralChange.IsZero() {
		return sdk.ErrInternal("invalid (empty) debt change")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}

	coins := sdk.NewCoin(msg.CollateralDenom, msg.CollateralChange)
	if !coins.IsValid() || coins.IsNegative() || coins.IsZero() {
		return sdk.ErrInternal("invalid (empty) coins")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgDepositCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawCollateral removes collateral to CSDT.
type MsgWithdrawCollateral struct {
	Sender           sdk.AccAddress `json:"sender" yaml:"sender"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	CollateralChange sdk.Int        `json:"collateral_change" yaml:"collateral_change"`
}

// NewMsgWithdrawCollateral returns a new MsgWithdrawCollateral.
func NewMsgWithdrawCollateral(sender sdk.AccAddress, collateralDenom string, collateralChange sdk.Int) MsgWithdrawCollateral {
	return MsgWithdrawCollateral{
		Sender:           sender,
		CollateralDenom:  collateralDenom,
		CollateralChange: collateralChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgWithdrawCollateral) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgWithdrawCollateral) Type() string { return "withdraw_collateral" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgWithdrawCollateral) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.CollateralChange.IsNegative() || msg.CollateralChange.IsZero() {
		return sdk.ErrInternal("invalid (empty) debt change")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}

	coins := sdk.NewCoin(msg.CollateralDenom, msg.CollateralChange)
	if !coins.IsValid() || coins.IsNegative() || coins.IsZero() {
		return sdk.ErrInternal("invalid (empty) coins")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgWithdrawCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgWithdrawCollateral) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgSettleDebt returns debt to CSDT.
type MsgSettleDebt struct {
	Sender          sdk.AccAddress `json:"sender" yaml:"sender"`
	CollateralDenom string         `json:"collateral_denom" yaml:"collateral_denom"`
	DebtDenom       string         `json:"debt_denom" yaml:"debt_denom"`
	DebtChange      sdk.Int        `json:"debt_change" yaml:"debt_change"`
}

// NewMsgSettleDebt returns a new MsgSettleDebt.
func NewMsgSettleDebt(sender sdk.AccAddress, collateralDenom, debtDenom string, debtChange sdk.Int) MsgSettleDebt {
	return MsgSettleDebt{
		Sender:          sender,
		CollateralDenom: collateralDenom,
		DebtDenom:       debtDenom,
		DebtChange:      debtChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgSettleDebt) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgSettleDebt) Type() string { return "settle_debt" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgSettleDebt) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.DebtChange.IsNegative() || msg.DebtChange.IsZero() {
		return sdk.ErrInternal("invalid (empty) debt change")
	}
	if len(msg.DebtDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) collateral denom")
	}

	coins := sdk.NewCoin(msg.DebtDenom, msg.DebtChange)
	if !coins.IsValid() || coins.IsNegative() || coins.IsZero() {
		return sdk.ErrInternal("invalid (empty) coins")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgSettleDebt) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgSettleDebt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawDebt withdraws debt from CSDT.
type MsgWithdrawDebt struct {
	Sender          sdk.AccAddress `json:"sender" yaml:"sender"`
	CollateralDenom string         `json:"collateral_denom" yaml:"collateral_denom"`
	DebtDenom       string         `json:"debt_denom" yaml:"debt_denom"`
	DebtChange      sdk.Int        `json:"debt_change" yaml:"debt_change"`
}

// NewMsgWithdrawDebt returns a new MsgWithdrawDebt.
func NewMsgWithdrawDebt(sender sdk.AccAddress, collateralDenom, debtDenom string, debtChange sdk.Int) MsgWithdrawDebt {
	return MsgWithdrawDebt{
		Sender:          sender,
		CollateralDenom: collateralDenom,
		DebtDenom:       debtDenom,
		DebtChange:      debtChange,
	}
}

// Route return the message type used for routing the message.
func (msg MsgWithdrawDebt) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgWithdrawDebt) Type() string { return "withdraw_debt" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgWithdrawDebt) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.DebtChange.IsNegative() || msg.DebtChange.IsZero() {
		return sdk.ErrInternal("invalid (empty) debt change")
	}
	if len(msg.DebtDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}

	coins := sdk.NewCoin(msg.DebtDenom, msg.DebtChange)
	if !coins.IsValid() || coins.IsNegative() || coins.IsZero() {
		return sdk.ErrInternal("invalid (empty) coins")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgWithdrawDebt) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgWithdrawDebt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgAddCollateralParam adds collateral to CSDT management
type MsgAddCollateralParam struct {
	Nominee          sdk.AccAddress `json:"nominee" yaml:"nominee"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	LiquidationRatio sdk.Dec        `json:"liquidation_ratio" yaml:"liquidation_ratio"`
	DebtLimit        sdk.Coins      `json:"debt_limit" yaml:"debt_limit"`
}

// NewMsgAddCollateralParam returns a new MsgAddCollateralParam.
func NewMsgAddCollateralParam(
	nominee sdk.AccAddress,
	collateralDenom string,
	liquidationRatio sdk.Dec,
	debtLimit sdk.Coins,
) MsgAddCollateralParam {
	return MsgAddCollateralParam{
		Nominee:          nominee,
		CollateralDenom:  collateralDenom,
		LiquidationRatio: liquidationRatio,
		DebtLimit:        debtLimit,
	}
}

// Route return the message type used for routing the message.
func (msg MsgAddCollateralParam) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgAddCollateralParam) Type() string { return "add_collateral_denom" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgAddCollateralParam) ValidateBasic() sdk.Error {
	if msg.Nominee.Empty() {
		return sdk.ErrInternal("invalid (empty) nominee address")
	}
	if msg.LiquidationRatio.IsNegative() || msg.LiquidationRatio.IsZero() {
		return sdk.ErrInternal("invalid (empty) liquidation ratio")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}

	if !msg.DebtLimit.IsValid() || msg.DebtLimit.IsAnyNegative() {
		return sdk.ErrInternal("invalid (empty) debt limit")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgAddCollateralParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgAddCollateralParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Nominee}
}

// MsgSetCollateralParam sets collateral in CSDT management
type MsgSetCollateralParam struct {
	Nominee          sdk.AccAddress `json:"nominee" yaml:"nominee"`
	CollateralDenom  string         `json:"collateral_denom" yaml:"collateral_denom"`
	LiquidationRatio sdk.Dec        `json:"liquidation_ratio" yaml:"liquidation_ratio"`
	DebtLimit        sdk.Coins      `json:"debt_limit" yaml:"debt_limit"`
}

// NewMsgSetCollateralParam returns a new MsgSetCollateralParam.
func NewMsgSetCollateralParam(
	nominee sdk.AccAddress,
	collateralDenom string,
	liquidationRatio sdk.Dec,
	debtLimit sdk.Coins,
) MsgSetCollateralParam {
	return MsgSetCollateralParam{
		Nominee:          nominee,
		CollateralDenom:  collateralDenom,
		LiquidationRatio: liquidationRatio,
		DebtLimit:        debtLimit,
	}
}

// Route return the message type used for routing the message.
func (msg MsgSetCollateralParam) Route() string { return ModuleName }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgSetCollateralParam) Type() string { return "set_collateral_denom" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgSetCollateralParam) ValidateBasic() sdk.Error {
	if msg.Nominee.Empty() {
		return sdk.ErrInternal("invalid (empty) nominee address")
	}
	if msg.LiquidationRatio.IsNegative() || msg.LiquidationRatio.IsZero() {
		return sdk.ErrInternal("invalid (empty) liquidation ratio")
	}
	if len(msg.CollateralDenom) == 0 {
		return sdk.ErrInternal("invalid (empty) debt denom")
	}

	if !msg.DebtLimit.IsValid() || msg.DebtLimit.IsAnyNegative() {
		return sdk.ErrInternal("invalid (empty) debt limit")
	}
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgSetCollateralParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgSetCollateralParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Nominee}
}

// MsgTransferCSDT changes the ownership of a csdt
type MsgTransferCSDT struct {
	// TODO
}
