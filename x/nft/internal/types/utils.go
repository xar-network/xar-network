/*

Copyright 2016 All in Bits, Inc
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

package types

import "strings"

// Findable is an interface for iterable types that allows the FindUtil function to work
type Findable interface {
	ElAtIndex(index int) string
	Len() int
}

// FindUtil is a binary search funcion for types that support the Findable interface (elements must be sorted)
func FindUtil(group Findable, el string) int {
	if group.Len() == 0 {
		return -1
	}
	low := 0
	high := group.Len() - 1
	median := 0
	for low <= high {
		median = (low + high) / 2
		switch compare := strings.Compare(group.ElAtIndex(median), el); {
		case compare == 0:
			// if group[median].element == el
			return median
		case compare == -1:
			// if group[median].element < el
			low = median + 1
		default:
			// if group[median].element > el
			high = median - 1
		}
	}
	return -1
}
