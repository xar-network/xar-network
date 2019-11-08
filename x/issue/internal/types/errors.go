package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Issue errors reserve 3500 ~ 3599.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeNotEnoughFee              sdk.CodeType = 3501
	CodeIssuerMismatch            sdk.CodeType = 3502
	CodeIssueIDNotValid           sdk.CodeType = 3503
	CodeIssueNameNotValid         sdk.CodeType = 3504
	CodeAmountNotValid            sdk.CodeType = 3505
	CodeIssueSymbolNotValid       sdk.CodeType = 3506
	CodeIssueTotalSupplyNotValid  sdk.CodeType = 3507
	CodeIssueDescriptionNotValid  sdk.CodeType = 3509
	CodeUnknownIssue              sdk.CodeType = 3510
	CanNotMint                    sdk.CodeType = 3511
	CanNotBurn                    sdk.CodeType = 3512
	CodeUnknownFeature            sdk.CodeType = 3513
	CodeUnknownFreezeType         sdk.CodeType = 3514
	CodeNotEnoughAmountToTransfer sdk.CodeType = 3515
	CodeCanNotFreeze              sdk.CodeType = 3516
	CodeFreezeEndTimeNotValid     sdk.CodeType = 3517
	CodeNotTransferIn             sdk.CodeType = 3518
	CodeNotTransferOut            sdk.CodeType = 3519
	CodeNegativeInflation         sdk.CodeType = 3520
)

//convert sdk.Error to error
func Errorf(err sdk.Error) error {
	return fmt.Errorf(err.Stacktrace().Error())
}

// Error constructors
func ErrOwnerMismatch() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssuerMismatch, "not the correct owner")
}
func ErrNotEnoughFee() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotEnoughFee, "insufficient funds for fees")
}
func ErrAmountNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAmountNotValid, "invalid amount")
}
func ErrCoinTotalSupplyMaxValueNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueTotalSupplyNotValid, "greater than total supply max value")
}
func ErrCoinSymbolNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueSymbolNotValid, "symbol length invalid")
}
func ErrCoinNamelNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueNameNotValid, "invalid coin name")
}
func ErrCoinDescriptionNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueDescriptionNotValid, "description is not valid json")
}
func ErrCoinDescriptionMaxLengthNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueDescriptionNotValid, "description over max length")
}
func ErrIssueID() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeIssueIDNotValid, "invalid issue ID")
}
func ErrCanNotMint() sdk.Error {
	return sdk.NewError(DefaultCodespace, CanNotMint, "can not mint coin")
}
func ErrCanNotBurn() sdk.Error {
	return sdk.NewError(DefaultCodespace, CanNotBurn, "can not burn coint")
}
func ErrUnknownIssue() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownIssue, "unknown issue ID")
}
func ErrUnknownFeatures() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownFeature, "unknown feature")
}
func ErrCanNotFreeze() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCanNotFreeze, "can not freeze the coin")
}
func ErrUnknownFreezeType() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownFreezeType, "unknown freeze type")
}
func ErrNotEnoughAmountToTransfer() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotEnoughAmountToTransfer, "insufficient balance to transfer")
}
func ErrCanNotTransferIn() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotTransferIn, "can not transfer in to account")
}
func ErrCanNotTransferOut() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotTransferOut, "can not transfer out of account")
}

func ErrNegativeInterest() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNegativeInflation, "cannot set negative interest")
}
