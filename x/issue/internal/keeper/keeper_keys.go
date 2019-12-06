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

package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

// Key for getting a the next available proposalID from the store
var (
	KeyDelimiter   = ":"
	KeyNextIssueID = []byte("newIssueID")
)

//func BytesString(b []byte) string {
//	return *(*string)(unsafe.Pointer(&b))
//}
// Key for getting a specific issuer from the store
func KeyIssuer(issueIdStr string) []byte {
	return []byte(fmt.Sprintf("issues:%s", issueIdStr))
}

// Key for getting a specific address from the store
func KeyAddressIssues(addr string) []byte {
	return []byte(fmt.Sprintf("address:%s", addr))
}

// Key for getting a specific allowed from the store
func KeyAllowed(issueID string, sender sdk.AccAddress, spender sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("allowed:%s:%s:%s", issueID, sender.String(), spender.String()))
}
func KeyFreeze(issueID string, accAddress sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("freeze:%s:%s", issueID, accAddress.String()))
}
func PrefixFreeze(issueID string) []byte {
	return []byte(fmt.Sprintf("freeze:%s", issueID))
}
func KeySymbolIssues(symbol string) []byte {
	return []byte(fmt.Sprintf("symbol:%s", strings.ToUpper(symbol)))
}

func KeyIssueIdStr(seq uint64) string {
	return fmt.Sprintf("%s%x", types.IDPreStr, seq)
}
