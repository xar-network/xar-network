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

type IteratorCB func(k []byte, v []byte) bool

type ArchiveStore interface {
	Get(key []byte) []byte
	Has(key []byte) bool
	Set(key []byte, value []byte)
	Delete(key []byte)
	Iterator(start []byte, end []byte, cb IteratorCB)
	ReverseIterator(start []byte, end []byte, cb IteratorCB)
	PrefixIterator(start []byte, cb IteratorCB)
	ReversePrefixIterator(start []byte, cb IteratorCB)
	Substore(prefix string) ArchiveStore
}
