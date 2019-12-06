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
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// implement fmt.Stringer
func (a PendingPriceAsset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s`, a.AssetCode))
}

// PendingPriceAsset struct that contains the info about the asset which price is still to be determined
type PendingPriceAsset struct {
	AssetCode string `json:"asset_code"`
}

func ValidateAddress(address string) (sdk.AccAddress, error) {
	oracle, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	return oracle, nil
}

func ParseOracles(addresses string) (Oracles, error) {
	res := make([]Oracle, 0)
	for _, address := range strings.Split(addresses, ",") {
		address = strings.TrimSpace(address)
		if len(address) == 0 {
			continue
		}
		oracleAddress, err := ValidateAddress(address)
		if err != nil {
			return nil, err
		}

		oracle := NewOracle(oracleAddress)

		res = append(res, oracle)
	}

	return res, nil
}
