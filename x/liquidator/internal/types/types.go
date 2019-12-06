/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SeizedDebt struct {
	Total         sdk.Int // Total debt seized from CSDTs. Known as Awe in maker.
	SentToAuction sdk.Int // Portion of seized debt that has had a (reverse) auction was started for it. Known as Ash in maker.
	// SentToAuction should always be < Total
}

// Available gets the seized debt that has not been sent for auction. Known as Woe in maker.
func (sd SeizedDebt) Available() sdk.Int {
	return sd.Total.Sub(sd.SentToAuction)
}

func (sd SeizedDebt) Settle(amount sdk.Int) (SeizedDebt, sdk.Error) {
	if amount.IsNegative() {
		return sd, sdk.ErrInternal("tried to settle a negative amount")
	}
	if amount.GT(sd.Total) {
		return sd, sdk.ErrInternal("tried to settle more debt than exists")
	}
	sd.Total = sd.Total.Sub(amount)
	sd.SentToAuction = sdk.MaxInt(sd.SentToAuction.Sub(amount), sdk.ZeroInt())
	return sd, nil
}
