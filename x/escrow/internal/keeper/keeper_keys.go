package keeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgard/hashgard/x/box/utils"

	"github.com/hashgard/hashgard/x/box/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Key for getting a the next available proposalID from the store
var (
	KeyDelimiter      = []byte(types.KeyDelimiterString)
	PrefixActiveQueue = []byte("active")
)

func KeyIdStr(boxType string, seq uint64) string {
	return fmt.Sprintf("%s%s%s", types.IDPreStr, types.GetMustBoxTypeValue(boxType), strconv.FormatUint(seq, 36))
}

// Key for getting a specific issuer from the store
func KeyNextBoxID(boxType string) []byte {
	return []byte(fmt.Sprintf("newBoxID:%s", boxType))
}
func KeyBox(boxIdStr string) []byte {
	return []byte(fmt.Sprintf("ids:%s:%s", utils.GetBoxTypeByValue(boxIdStr), boxIdStr))
}
func KeyAllBox() []byte {
	return []byte("ids:")
}

// Key for getting a specific address from the store
func KeyAddress(boxType string, accAddress sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("address:%s:%s", boxType, accAddress.String()))
}
func KeyName(boxType string, name string) []byte {
	return []byte(fmt.Sprintf("name:%s:%s", boxType, strings.ToLower(name)))
}

//func KeyAddressInject(id string, accAddress sdk.AccAddress) []byte {
//	return []byte(fmt.Sprintf("deposit:%s:%s", id, accAddress.String()))
//}
//
//func GetAddressFromKeyAddressInject(keyAddressInject []byte) sdk.AccAddress {
//	str := fmt.Sprintf("%s", keyAddressInject)
//	keys := strings.Split(str, ":")
//	address, _ := sdk.AccAddressFromBech32(keys[2])
//	return address
//}
//func PrefixKeyDeposit(id string) []byte {
//	return []byte(fmt.Sprintf("deposit:%s", id))
//}

// Returns the key for a id in the activeQueue
func PrefixActiveBoxQueueTime(endTime int64) []byte {
	return []byte(fmt.Sprintf("active:%d", endTime))
}

// Returns the key for a proposalID in the activeQueue
func KeyActiveBoxQueue(endTime int64, boxIdStr string) []byte {
	return []byte(fmt.Sprintf("active:%d:%s", endTime, boxIdStr))
}
