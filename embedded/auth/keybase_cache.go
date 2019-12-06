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
	"errors"
	"net/http"
	"sync"

	"github.com/tendermint/tendermint/crypto"

	"github.com/xar-network/xar-network/embedded/session"
)

var kb *Keybase
var currID string
var mtx sync.RWMutex

func GetKBFromSession(r *http.Request) (*Keybase, error) {
	id, err := session.GetStr(r, keybaseIDKey)
	if err != nil {
		return nil, err
	}
	kb := GetKB(id)
	if kb == nil {
		return nil, errors.New("no keybase found")
	}
	return kb, nil
}

func MustGetKBFromSession(r *http.Request) *Keybase {
	kb, err := GetKBFromSession(r)
	if err != nil {
		panic(err)
	}
	return kb
}

func MustGetKBPassphraseFromSession(r *http.Request) string {
	return session.MustGetStr(r, keybasePassphraseKey)
}

func GetKB(id string) *Keybase {
	mtx.RLock()
	defer mtx.RUnlock()
	if currID != id {
		return nil
	}

	return kb
}

func ReplaceKB(name string, passphrase string, pk crypto.PrivKey) string {
	mtx.Lock()
	defer mtx.Unlock()
	currID = ReadStr32()
	kb = NewHotKeybase(name, passphrase, pk)
	return currID
}
