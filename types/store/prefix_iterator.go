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

package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func KVStorePrefixIterator(kvs KVStore, prefix []byte) sdk.Iterator {
	return kvs.Iterator(prefix, sdk.PrefixEndBytes(prefix))
}

func KVStoreReversePrefixIterator(kvs KVStore, prefix []byte) sdk.Iterator {
	return kvs.ReverseIterator(prefix, sdk.PrefixEndBytes(prefix))
}
