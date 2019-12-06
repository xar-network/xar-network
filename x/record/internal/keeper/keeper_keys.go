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

	"github.com/xar-network/xar-network/x/record/internal/types"
)

// Key for getting next available recordID from the store
var (
	KeyDelimiter    = ":"
	KeyNextRecordID = []byte("newRecordID")
)

// Key for getting a specific record content hash from the store
func KeyRecord(recordHash string) []byte {
	return []byte(fmt.Sprintf("hash:%s", recordHash))
}

// Key for getting records by a specific address from the store
func KeyAddress(addr string) []byte {
	return []byte(fmt.Sprintf("addr:%s", addr))
}

// Key for saving a record by a specific address:id
func KeyAddressRecord(addr string, id string) []byte {
	return []byte(fmt.Sprintf("addr:%s:%s", addr, id))
}

func KeyRecordId(id string) []byte {
	return []byte(fmt.Sprintf("id:%s", id))
}

func KeyRecordIdStr(seq uint64) string {
	return fmt.Sprintf("%s%x", types.IDPreStr, seq)
}
