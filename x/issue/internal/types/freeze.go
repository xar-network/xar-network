package types

import (
	"fmt"
	"strings"
)

const (
	FreezeIn       = "in"
	FreezeOut      = "out"
	FreezeInAndOut = "in-out"
)

var FreezeTypes = map[string]int{FreezeIn: 1, FreezeOut: 1, FreezeInAndOut: 1}

type IssueAddressFreeze struct {
	Address string `json:"address"`
}

type IssueAddressFreezeList []IssueAddressFreeze

func (ci IssueAddressFreeze) String() string {
	return fmt.Sprintf(`FreezeList:\n
	Address:			%s`,
		ci.Address)
}

//nolint
func (ci IssueAddressFreezeList) String() string {
	out := fmt.Sprintf("%-44s\n",
		"Address")
	for _, v := range ci {
		out += fmt.Sprintf("%-44s\n",
			v.Address)
	}
	return strings.TrimSpace(out)
}
