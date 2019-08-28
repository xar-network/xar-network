package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/issue/internal/types"
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
