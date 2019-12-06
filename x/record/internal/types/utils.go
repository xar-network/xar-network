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
