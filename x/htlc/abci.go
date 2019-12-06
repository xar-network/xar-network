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

package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker handles block beginning logic
func BeginBlocker(ctx sdk.Context, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/htlc"))

	currentBlockHeight := uint64(ctx.BlockHeight())
	iterator := k.IterateHTLCExpireQueueByHeight(ctx, currentBlockHeight)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// get the hash lock
		var hashLock []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &hashLock)

		htlc, _ := k.GetHTLC(ctx, hashLock)

		// update the state
		htlc.State = EXPIRED
		k.SetHTLC(ctx, htlc, hashLock)

		// delete from the expiration queue
		k.DeleteHTLCFromExpireQueue(ctx, currentBlockHeight, hashLock)

		// add tags
		tags = tags.AppendTags(sdk.NewTags(
			TagHashLock, []byte(hex.EncodeToString(hashLock)),
		))

		ctx.Logger().Info(fmt.Sprintf("HTLC [%s] is expired", hex.EncodeToString(hashLock)))
	}

	return
}
