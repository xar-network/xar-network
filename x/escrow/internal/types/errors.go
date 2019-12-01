package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeNotEnoughFee           sdk.CodeType = 1
	CodeBoxOwnerMismatch       sdk.CodeType = 2
	CodeBoxIDNotValid          sdk.CodeType = 3
	CodeBoxNameNotValid        sdk.CodeType = 4
	CodeAmountNotValid         sdk.CodeType = 5
	CodeDecimalsNotValid       sdk.CodeType = 6
	CodeTimelineNotValid       sdk.CodeType = 7
	CodeBoxDescriptionNotValid sdk.CodeType = 8
	CodeUnknownBox             sdk.CodeType = 9
	CodeUnknownBoxType         sdk.CodeType = 10
	CodeUnknownOperation       sdk.CodeType = 11
	CodeInterestInjectNotValid sdk.CodeType = 12
	CodeInterestCancelNotValid sdk.CodeType = 13
	CodeNotEnoughAmount        sdk.CodeType = 14
	CodeTimeNotValid           sdk.CodeType = 15
	CodeNotAllowedOperation    sdk.CodeType = 16
	CodeNotSupportOperation    sdk.CodeType = 17
	CodeUnknownFeature         sdk.CodeType = 18
	CodeNotTransfer            sdk.CodeType = 19
)

//convert sdk.Error to error
func Errorf(err sdk.Error) error {
	return fmt.Errorf(err.Stacktrace().Error())
}

// Error constructors
func ErrOwnerMismatch(id string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxOwnerMismatch, fmt.Sprintf("Owner mismatch with box %s", id))
}
func ErrNotEnoughFee() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotEnoughFee, fmt.Sprintf("Not enough fee"))
}
func ErrDecimalsNotValid(decimals uint) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDecimalsNotValid, "%d is not a valid decimals", decimals)
}
func ErrTimelineNotValid(time []int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTimelineNotValid, "%d is not a valid time line", time)
}
func ErrTimeNotValid(timeKey string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTimeNotValid, "%s is not a valid timestamp", timeKey)
}
func ErrAmountNotValid(key string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAmountNotValid, "%s is not a valid amount", key)
}
func ErrInterestInjectNotValid(coin sdk.Coin) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInterestInjectNotValid, "%s is not a valid interest injection", coin.String())
}
func ErrInterestCancelNotValid(coin sdk.Coin) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInterestCancelNotValid, "%s is not a valid interest fetch", coin.String())
}
func ErrBoxPriceNotValid(name string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxNameNotValid, fmt.Sprintf("Price mismatch with %s", name))
}
func ErrBoxNameNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxNameNotValid, fmt.Sprintf("Name max length is %d", BoxNameMaxLength))
}
func ErrBoxDescriptionNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxDescriptionNotValid, "Description is not valid json")
}
func ErrBoxDescriptionMaxLengthNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxDescriptionNotValid, "Description max length is %d", BoxDescriptionMaxLength)
}
func ErrBoxID(id string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBoxIDNotValid, fmt.Sprintf("id %s is not a valid id", id))
}
func ErrUnknownBox(id string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownBox, fmt.Sprintf("Unknown box with id %s", id))
}
func ErrUnknownBoxType() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownBoxType, fmt.Sprintf("Unknown type"))
}
func ErrUnknownOperation() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownOperation, fmt.Sprintf("Unknown operation"))
}
func ErrNotEnoughAmount() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotEnoughAmount, fmt.Sprintf("Not enough amount"))
}

func ErrNotSupportOperation() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotSupportOperation, fmt.Sprintf("Not support operation"))
}
func ErrNotAllowedOperation(status string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotAllowedOperation, fmt.Sprintf("Not allowed operation in current status: %s", status))
}
func ErrUnknownFeatures() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownFeature, fmt.Sprintf("Unknown feature"))
}
func ErrCanNotTransfer(id string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotTransfer, fmt.Sprintf("The box %s Can not be transfer", id))
}
