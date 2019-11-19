package types

import (
	"bytes"

	"github.com/olekukonko/tablewriter"

	"github.com/xar-network/xar-network/types/store"
)

type NamedMarket struct {
	ID              store.EntityID
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
			m.ID.String(),
			m.Name,
			m.BaseAssetDenom,
			m.QuoteAssetDenom,
		})
	}

	t.Render()
	return string(buf.Bytes())
}
