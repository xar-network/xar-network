package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace codespace for the module
	DefaultCodespace sdk.CodespaceType = ModuleName

	// CodeEmptyInput error code for empty input errors
	CodeNoPOA sdk.CodeType = 1

	CodeTokenSymbolDoesNotExist sdk.CodeType = 101
)

// ErrEmptyInput Error constructor
func ErrNoPOA(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoPOA, fmt.Sprintf("Invalid POA address."))
}

func ErrTokenSymbolDoesNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenSymbolDoesNotExist, "Token symbol does not exist")
}
