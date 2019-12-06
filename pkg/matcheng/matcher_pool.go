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

package matcheng

import "sync"

var pool *sync.Pool

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return NewMatcher()
		},
	}
}

func GetMatcher() *Matcher {
	return pool.Get().(*Matcher)
}

func ReturnMatcher(m *Matcher) {
	m.Reset()
	pool.Put(m)
}
