/*

Copyright 2019 All in Bits, Inc
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

package matcheng

import (
	"math/rand"
	"testing"

	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BenchmarkMatching(b *testing.B) {
	id := store.NewEntityID(0)
	matcher := GetMatcher()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		matcher.Reset()
		for j := 0; j < 10000; j++ {
			id = id.Inc()
			matcher.EnqueueOrder(Bid, id, sdk.NewUint(uint64(j)), sdk.NewUint(uint64(j)))
		}
		for j := 100; j < 11000; j++ {
			id := id.Inc()
			matcher.EnqueueOrder(Ask, id, sdk.NewUint(uint64(j)), sdk.NewUint(uint64(j)))
		}
		b.StartTimer()
		matcher.Match()
	}
}

func BenchmarkQueueing(b *testing.B) {
	id := store.NewEntityID(0)
	matcher := GetMatcher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		matcher.Reset()
		b.StartTimer()
		for j := 0; j < 100; j++ {
			id = id.Inc()
			price := sdk.NewUint(rand.Uint64())
			quantity := sdk.NewUint(rand.Uint64())
			matcher.EnqueueOrder(Bid, id.Inc(), price, quantity)
		}
	}
}
