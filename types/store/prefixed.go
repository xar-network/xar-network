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

import sdk "github.com/cosmos/cosmos-sdk/types"

type Prefixed struct {
	backend KVStore
	prefix  []byte
}

func NewPrefixed(backend KVStore, prefix []byte) *Prefixed {
	return &Prefixed{
		backend: backend,
		prefix:  prefix,
	}
}

func (p *Prefixed) Get(key []byte) []byte {
	return p.backend.Get(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Has(key []byte) bool {
	return p.backend.Has(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Set(key, value []byte) {
	p.backend.Set(PrefixKeyBytes(p.prefix, key), value)
}

func (p *Prefixed) Delete(key []byte) {
	p.backend.Delete(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Iterator(start, end []byte) sdk.Iterator {
	return p.backend.Iterator(PrefixKeyBytes(p.prefix, start), PrefixKeyBytes(p.prefix, end))
}

func (p *Prefixed) ReverseIterator(start, end []byte) sdk.Iterator {
	return p.backend.ReverseIterator(PrefixKeyBytes(p.prefix, start), PrefixKeyBytes(p.prefix, end))
}
