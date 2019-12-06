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

package conv

import (
	"encoding/binary"
	"io"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SDKUint2Big(in sdk.Uint) *big.Int {
	out, _ := new(big.Int).SetString(in.String(), 10)
	return out
}

func Uint642Bytes(in uint64) []byte {
	b := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, in)
	return b
}

func ReadUint64(r io.Reader) (uint64, error) {
	b := make([]byte, 8, 8)
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}
