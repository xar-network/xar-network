package types

import (
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
func IsRecordId(id string) bool {
	return strings.HasPrefix(id, IDPreStr)
}

func CheckRecordId(issueID string) sdk.Error {
	if !IsRecordId(issueID) {
		return ErrRecordIDNotValid(issueID)
	}
	return nil
}

func CheckRecordHash(hash string) sdk.Error {
	if len(hash) != 64 {
		return ErrRecordHashNotValid()
	}
	return nil
}
func GetRecordTags(info *RecordInfo) sdk.Events {
	res := sdk.NewEvent(
		"tags",
		sdk.NewAttribute("tags.Sender", info.Sender.String()),
		sdk.NewAttribute("tags.ID", info.ID),
		sdk.NewAttribute("tags.Hash", info.Hash),
	)
	if len(info.RecordType) > 0 {
		res.AppendAttributes(sdk.NewAttribute("tags.RecordType", info.RecordType))
	}
	if len(info.RecordNo) > 0 {
		res.AppendAttributes(sdk.NewAttribute("tags.RecordNo", info.RecordNo))
	}
	if len(info.Author) > 0 {
		res.AppendAttributes(sdk.NewAttribute("tags.Author", info.Author))
	}
	return sdk.Events{res}
}
