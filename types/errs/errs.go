/*

Copyright 2019 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package errs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	_ sdk.CodeType = iota
	CodeNotFound
	CodeAlreadyExists
	CodeInvalidArgument
	CodeMarshalFailure
	CodeUnmarshalFailure

	CodespaceUEX sdk.CodespaceType = "dex-demo"
)

func newErrWithUEXCodespace(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(CodespaceUEX, code, msg)
}

func ErrNotFound(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeNotFound, msg)
}

func ErrAlreadyExists(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeAlreadyExists, msg)
}

func ErrInvalidArgument(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeInvalidArgument, msg)
}

func ErrMarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeMarshalFailure, msg)
}

func ErrUnmarshalFailure(msg string) sdk.Error {
	return newErrWithUEXCodespace(CodeUnmarshalFailure, msg)
}

func ErrOrBlankResult(err sdk.Error) sdk.Result {
	if err == nil {
		return sdk.Result{}
	}

	return err.Result()
}
