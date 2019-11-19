package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/embedded/auth"
)

func RandAddr() sdk.AccAddress {
	return sdk.AccAddress(auth.ReadN(sdk.AddrLen))
}

func Rand32() []byte {
	return auth.ReadN(32)
}
