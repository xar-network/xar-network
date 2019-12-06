/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

// GetUniId returns the unique uni id for the provided denominations.
// The uni id is in the format of 'u-coin-name' which the denomination
// is not iris-atto.
func GetUniId(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denomnations for forming uni id are equal")
	}

	denom := denom1
	if denom == "uftm" {
		denom = denom2
	}
	return fmt.Sprintf("u-%s", denom), nil
}

// GetCoinMinDenomFromUniDenom returns the token denom by uni denom
func GetCoinMinDenomFromUniDenom(uniDenom string) (string, sdk.Error) {
	err := CheckUniDenom(uniDenom)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(uniDenom, "u-"), nil
}

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) sdk.Error {
	return nil
}

// CheckUniId returns nil if the uni id is valid
func CheckUniId(uniId string) sdk.Error {
	return nil
}

// GetUniDenom returns uni denom if the uni id is valid
func GetUniDenom(uniId string) (string, sdk.Error) {
	if err := CheckUniId(uniId); err != nil {
		return "", err
	}

	uniDenom := fmt.Sprintf("u-%s", uniId)
	return uniDenom, nil
}
