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
	"bytes"
	"encoding/binary"
)

func PrefixKeyString(prefix string, subkeys ...[]byte) []byte {
	buf := [][]byte{[]byte(prefix)}
	return PrefixKeyBytes(append(buf, subkeys...)...)
}

func PrefixKeyBytes(subkeys ...[]byte) []byte {
	if len(subkeys) == 0 {
		return []byte{}
	}

	var buf bytes.Buffer
	buf.Write(subkeys[0])

	if len(subkeys) > 1 {
		for _, sk := range subkeys[1:] {
			if len(sk) == 0 {
				continue
			}

			buf.WriteRune('/')
			buf.Write(sk)
		}
	}

	return buf.Bytes()
}

func IntSubkey(subkey int) []byte {
	if subkey < 0 {
		panic("cannot use negative numbers in subkeys")
	}
	return Uint64Subkey(uint64(subkey))
}

func Int64Subkey(subkey int64) []byte {
	if subkey < 0 {
		panic("cannot use negative numbers in subkeys")
	}
	return Uint64Subkey(uint64(subkey))
}

func Uint64Subkey(subkey uint64) []byte {
	b := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, subkey)
	return b
}
