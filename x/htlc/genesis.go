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

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for hashLockHex, htlc := range data.PendingHTLCs {
		hashLock, err := hex.DecodeString(hashLockHex)
		if err != nil {
			panic(fmt.Errorf("failed to initialize HTLC genesis state: %s", err.Error()))
		}

		k.SetHTLC(ctx, htlc, hashLock)
		k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingHTLCs := make(map[string]HTLC)

	k.IterateHTLCs(ctx, func(hlock []byte, h HTLC) (stop bool) {
		if h.State == OPEN || h.State == EXPIRED {
			if h.State == OPEN {
				h.ExpireHeight = h.ExpireHeight - uint64(ctx.BlockHeight()) + 1
				pendingHTLCs[hex.EncodeToString(hlock)] = h
			} else {
				_, err := k.RefundHTLC(ctx, hlock)
				if err != nil {
					panic(fmt.Errorf("failed to export HTLC genesis state: %s", hex.EncodeToString(hlock)))
				}
			}
		}

		return false
	})

	return GenesisState{
		PendingHTLCs: pendingHTLCs,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingHTLCs: map[string]HTLC{},
	}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		PendingHTLCs: map[string]HTLC{},
	}
}
