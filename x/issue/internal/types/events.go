/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

// Issue module event types
var (
	TxCategory     = "issue"
	EventTypeIssue = "slash"

	Action          = "action"
	Category        = "category"
	Sender          = "sender"
	Owner           = "owner"
	Fee             = "fee"
	IssueID         = "issue-id"
	Feature         = "feature"
	Name            = "name"
	Symbol          = "symbol"
	TotalSupply     = "total-supply"
	MintingFinished = "minting-finished"
	FreezeType      = "freeze-type"

	AttributeValueCategory = ModuleName
)
