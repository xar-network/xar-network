package keeper

import "github.com/cosmos/cosmos-sdk/types"

type AddLiquidity interface {
	GetSenderAddress() types.AccAddress
	GetDeposit() types.Coin
	GetNativeAmount() types.Int
}
