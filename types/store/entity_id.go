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
	"math/big"

	"github.com/xar-network/xar-network/pkg/conv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ZeroEntityID = NewEntityID(0)

type EntityID sdk.Uint

func NewEntityID(id uint64) EntityID {
	return EntityID(sdk.NewUint(id))
}

func NewEntityIDFromString(str string) EntityID {
	return EntityID(sdk.NewUintFromString(str))
}

func NewEntityIDFromBytes(b []byte) EntityID {
	s := new(big.Int).SetBytes(b)
	return EntityID(sdk.NewUintFromBigInt(s))
}

func (id EntityID) String() string {
	return sdk.Uint(id).String()
}

// Uint64 converts Uint to uint64
// Panics if the value is out of range
func (id EntityID) Uint64() uint64 {
	return sdk.Uint(id).Uint64()
}

func (id EntityID) Bytes() []byte {
	var buf [32]byte
	bn := conv.SDKUint2Big(sdk.Uint(id))
	b := bn.Bytes()
	copy(buf[32-len(b):], b)
	return buf[:]
}

func (id EntityID) Inc() EntityID {
	return EntityID(sdk.Uint(id).Add(sdk.OneUint()))
}

func (id EntityID) Dec() EntityID {
	if !id.IsDefined() {
		return id
	}

	return EntityID(sdk.Uint(id).Sub(sdk.OneUint()))
}

func (id EntityID) Cmp(b EntityID) int {
	uintA := sdk.Uint(id)
	uintB := sdk.Uint(b)

	if uintA.GT(uintB) {
		return 1
	}

	if uintA.LT(uintB) {
		return -1
	}

	return 0
}

func (id EntityID) IsDefined() bool {
	return !sdk.Uint(id).IsZero()
}

func (id EntityID) IsZero() bool {
	return sdk.Uint(id).IsZero()
}

func (id EntityID) Equals(other EntityID) bool {
	return sdk.Uint(id).Equal(sdk.Uint(other))
}

func (id EntityID) MarshalAmino() (string, error) {
	return sdk.Uint(id).MarshalAmino()
}

func (id *EntityID) UnmarshalAmino(text string) error {
	var u sdk.Uint
	err := u.UnmarshalAmino(text)
	if err != nil {
		return err
	}

	*id = EntityID(u)
	return nil
}

func (id *EntityID) UnmarshalJSON(data []byte) error {
	return (*sdk.Uint)(id).UnmarshalJSON(data)
}

func (id EntityID) MarshalJSON() ([]byte, error) {
	return sdk.Uint(id).MarshalJSON()
}
