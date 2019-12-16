package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	IncorrectBaseAmountForFeeCode = 101 + iota
)

const (
	MsgIncorrectBaseAmountForFee = "base fee amount cannot be less or equal to zero"
	MsgIncorrectMinimumFee = "minimum fee cannot be less than zero"
)

var ErrIncorrectBaseAmountForFee = sdk.NewError(DefaultCodespace, IncorrectBaseAmountForFeeCode, MsgIncorrectBaseAmountForFee)
