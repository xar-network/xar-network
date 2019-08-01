package types

import (
	"math/big"
	"math/rand"
	"strings"

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
	return strings.HasPrefix(issueID, IDPreStr)
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
