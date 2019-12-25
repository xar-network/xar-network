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
See the License for the specific language governing Permissions and
limitations under the License.

*/

package types

import (
	"bytes"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type ListQueryResult struct {
	Orders []Order `json:"orders"`
}

func (l ListQueryResult) String() string {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	t.SetHeader([]string{
		"ID",
		"Owner",
		"MarketID",
		"Direction",
		"Price",
		"Quantity",
		"Time In Force",
		"Created Block",
	})

	for _, o := range l.Orders {
		t.Append([]string{
			o.ID.String(),
			o.Owner.String(),
			o.MarketID.String(),
			o.Direction.String(),
			o.Price.String(),
			o.Quantity.String(),
			strconv.FormatUint(uint64(o.TimeInForceBlocks), 10),
			strconv.Itoa(int(o.CreatedBlock)),
		})
	}
	t.Render()
	return string(buf.Bytes())
}
