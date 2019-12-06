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
	"encoding/hex"
	"strings"
)

func HexToBytes(in string) ([]byte, error) {
	if strings.HasPrefix(in, "0x") {
		in = in[2:]
	}

	return hex.DecodeString(in)
}

func MustHexToBytes(in string) []byte {
	out, err := HexToBytes(in)
	if err != nil {
		panic(err)
	}
	return out
}
