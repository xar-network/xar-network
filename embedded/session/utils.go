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

package session

import (
	"errors"
	"fmt"
	"net/http"
)

func MustGetStr(r *http.Request, key string) string {
	out, err := GetStr(r, key)
	if err != nil {
		panic(err)
	}
	return out
}

func GetStr(r *http.Request, key string) (string, error) {
	store, _ := SessionStore.Get(r, sessionName)
	val, ok := store.Values[key]
	if !ok || val == "" {
		return "", errors.New(fmt.Sprintf("key %s not found in session", key))
	}
	return val.(string), nil
}

func SetStrings(w http.ResponseWriter, r *http.Request, kvPairs ...string) error {
	if len(kvPairs) < 2 || len(kvPairs)%2 != 0 {
		return errors.New("mismatched KV pairs")
	}

	store, _ := SessionStore.Get(r, sessionName)
	for i := 0; i < len(kvPairs); i += 2 {
		k := kvPairs[i]
		v := kvPairs[i+1]
		store.Values[k] = v
	}

	return store.Save(r, w)
}
