package types

import (
	"math/big"
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

func GetDecimalsInt(decimals uint) sdk.Int {
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return sdk.NewIntFromBigInt(exp)
}

func MulDecimals(totalSupply sdk.Int, decimals uint) sdk.Int {
	return totalSupply.Mul(GetDecimalsInt(decimals))
}

func QuoDecimals(totalSupply sdk.Int, decimals uint) sdk.Int {
	return totalSupply.Quo(GetDecimalsInt(decimals))
}

func IssueOwnerCheck(cliCtx context.CLIContext, sender sdk.AccAddress, issueID string) (Issue, error) {
	var issueInfo Issue
	// Query the issue
	res, _, err := cliCtx.QueryWithData(GetQueryIssuePath(issueID), nil)
	if err != nil {
		return nil, err
	}
	cliCtx.Codec.MustUnmarshalJSON(res, &issueInfo)

	if !sender.Equals(issueInfo.GetOwner()) {
		return nil, Errorf(ErrOwnerMismatch())
	}
	return issueInfo, nil
}
