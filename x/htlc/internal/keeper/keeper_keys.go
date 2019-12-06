/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter          = []byte(":")                // key separator
	PrefixHTLC            = []byte("htlcs:")           // key prefix for HTLC
	PrefixHTLCExpireQueue = []byte("htlcExpireQueue:") // key prefix for the HTLC expiration queue
)

// KeyHTLC returns the key for an HTLC by the specified hash lock
func KeyHTLC(hashLock []byte) []byte {
	return append(PrefixHTLC, hashLock...)
}

// KeyHTLCExpireQueue returns the key for HTLC expiration queue by the specified height and hash lock
func KeyHTLCExpireQueue(expireHeight uint64, hashLock []byte) []byte {
	prefix := append(PrefixHTLCExpireQueue, sdk.Uint64ToBigEndian(expireHeight)...)
	return append(append(prefix, KeyDelimiter...), hashLock...)
}

// KeyHTLCExpireQueueSubspace returns the key prefix for HTLC expiration queue by the given height
func KeyHTLCExpireQueueSubspace(expireHeight uint64) []byte {
	return append(append(PrefixHTLCExpireQueue, sdk.Uint64ToBigEndian(expireHeight)...), KeyDelimiter...)
}
