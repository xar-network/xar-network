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

import (
	"bytes"
	"fmt"
)

func PlotCurves(bids []AggregatePrice, asks []AggregatePrice) string {
	var buf bytes.Buffer
	buf.WriteString("\"Ask\"\n")

	for i, entry := range asks {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("%s 0\n", entry[0]))
		}
		if i > 0 {
			buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], asks[i-1][1]))
		}
		buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], entry[1]))
	}

	buf.WriteString("\n\n")
	buf.WriteString("\"Bid\"\n")

	for i := len(bids) - 1; i >= 0; i-- {
		entry := bids[i]
		if i == len(bids)-1 {
			buf.WriteString(fmt.Sprintf("%s 0\n", entry[0]))
		}
		if i != len(bids)-1 {
			buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], bids[i+1][1]))
		}
		buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], entry[1]))
		if i == 0 {
			buf.WriteString(fmt.Sprintf("0 %s\n", entry[1]))
		}
	}

	out := buf.Bytes()
	return string(out)
}
