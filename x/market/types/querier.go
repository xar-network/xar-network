package types

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type NamedMarket struct {
	ID              string
	BaseAssetDenom  string
	QuoteAssetDenom string
	Name            string
}

type ListQueryResult struct {
	Markets []NamedMarket `json:"markets"`
}

func (l ListQueryResult) String() string {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	t.SetHeader([]string{
		"ID",
		"Name",
		"Base Asset ID",
		"Quote Asset ID",
	})

	for _, m := range l.Markets {
		t.Append([]string{
			m.ID,
			m.Name,
			m.BaseAssetDenom,
			m.QuoteAssetDenom,
		})
	}

	t.Render()
	return string(buf.Bytes())
}
