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

package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func ReadStr32() string {
	return ReadStrN(32)
}

func ReadStrN(byteLen int) string {
	return hex.EncodeToString(ReadN(byteLen))
}

func ReadN(n int) []byte {
	buf := make([]byte, n, n)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}
