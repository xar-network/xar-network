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

type KVStore interface {
	Get(key []byte) []byte

	Has(key []byte) bool

	Set(key, value []byte)

	Delete(key []byte)

	Iterator(start, end []byte) sdk.Iterator

	ReverseIterator(start, end []byte) sdk.Iterator
}
