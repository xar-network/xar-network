package fee

const (
	MsgIncorrectBaseAmountForFee = "base fee amount cannot be less or equal to zero"
	MsgIncorrectMinimumFee       = "minimum fee cannot be less than zero"
	MsgNumeratorLTEDenominator   = "fee numerator cannot be less or equal to fee denominator"
	MsgAmountSubFeeTooSmall      = "amount sub fee is too small"
	MsgRatioIsIncorrect          = "incorrect fee ratio"
)
