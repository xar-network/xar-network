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
