/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

package types

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	randomBytes = []rune("abcdefghijklmnopqrstuvwxyz")
)

func GetRandomString(l int) string {
	result := make([]rune, l)
	length := len(randomBytes)
	for i := range result {
		result[i] = randomBytes[rand.Intn(length)]
	}
	return string(result)
}
func IsIssueId(issueID string) bool {
	return true //strings.HasPrefix(issueID, IDPreStr)
}

func CheckIssueId(issueID string) sdk.Error {
	if !IsIssueId(issueID) {
		return ErrIssueID()
	}
	return nil
}

func IssueOwnerCheck(cliCtx context.CLIContext, sender sdk.AccAddress, issueID string) (Issue, error) {
	var issueInfo Issue
	// Query the issue
	res, height, err := cliCtx.QueryWithData(GetQueryIssuePath(issueID), nil)
	if err != nil {
		return nil, err
	}
	cliCtx = cliCtx.WithHeight(height)

	cliCtx.Codec.MustUnmarshalJSON(res, &issueInfo)

	if !sender.Equals(issueInfo.GetOwner()) {
		return nil, Errorf(ErrOwnerMismatch())
	}
	return issueInfo, nil
}
