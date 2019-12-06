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

package serde

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HexBytes []byte

func (h *HexBytes) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		*h = nil
		return nil
	}
	if len(buf) <= 2 {
		return errors.New("no value")
	}

	unquoted := string(buf[1 : len(buf)-1])
	data, err := hex.DecodeString(unquoted[2:])
	if err != nil {
		return err
	}
	*h = data
	return nil
}

func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(fmt.Sprintf("0x%s", hex.EncodeToString(h)))
}

func MustMarshalSortedJSON(in interface{}) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
