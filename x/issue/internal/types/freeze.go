package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	FreezeIn       = "in"
	FreezeOut      = "out"
	FreezeInAndOut = "in-out"

	UnFreezeEndTime int64 = 0
)

var FreezeTypes = map[string]int{FreezeIn: 1, FreezeOut: 1, FreezeInAndOut: 1}

type IssueFreeze struct {
	OutEndTime int64 `json:"out_end_time"`
	InEndTime  int64 `json:"in_end_time"`
}
type IssueAddressFreeze struct {
	Address    string `json:"address"`
	OutEndTime int64  `json:"out_end_time"`
	InEndTime  int64  `json:"in_end_time"`
}

type IssueAddressFreezeList []IssueAddressFreeze

func NewIssueFreeze(outEndTime int64, inEndTime int64) IssueFreeze {
	return IssueFreeze{outEndTime, inEndTime}
}

func (ci IssueFreeze) String() string {
	return fmt.Sprintf(`Freeze:\n
	Out-end-time:			%T
	In-end-time:			%T`,
		time.Unix(ci.OutEndTime, 0), time.Unix(ci.InEndTime, 0))
}
func (ci IssueAddressFreeze) String() string {
	return fmt.Sprintf(`FreezeList:\n
	Address:			%s
	Out-end-time:			%T
	In-end-time:			%T`,
		ci.Address, time.Unix(ci.OutEndTime, 0), time.Unix(ci.InEndTime, 0))
}

//nolint
func (ci IssueAddressFreezeList) String() string {
	out := fmt.Sprintf("%-44s|%-32s|%-32s\n",
		"Address", "Out-end-time", "In-end-time")
	for _, v := range ci {
		out += fmt.Sprintf("%-44s|%-32s|%-32s\n",
			v.Address, time.Unix(v.OutEndTime, 0), time.Unix(v.InEndTime, 0))
	}
	return strings.TrimSpace(out)
}
