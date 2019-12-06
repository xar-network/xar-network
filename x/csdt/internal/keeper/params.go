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

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// ---------- Module Parameters ----------
// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramsSubspace.GetParamSet(ctx, &p)
	return p
}

// SetParams sets params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramsSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) AddCollateralParam(ctx sdk.Context, nominee string, collateralParam types.CollateralParam) sdk.Error {
	if !k.IsNominee(ctx, nominee) {
		return sdk.ErrInternal(fmt.Sprintf("not a nominee: '%s'", nominee))
	}
	params := k.GetParams(ctx)
	if params.IsCollateralPresent(collateralParam.Denom) {
		return sdk.ErrInternal(fmt.Sprintf("param already exists: '%s'", collateralParam.String()))
	}
	params.CollateralParams = append(params.CollateralParams, collateralParam)
	k.SetParams(ctx, params)
	return nil
}

func (k Keeper) SetCollateralParam(ctx sdk.Context, nominee string, collateralParam types.CollateralParam) sdk.Error {
	if !k.IsNominee(ctx, nominee) {
		return sdk.ErrInternal(fmt.Sprintf("not a nominee: '%s'", nominee))
	}
	params := k.GetParams(ctx)
	if !params.IsCollateralPresent(collateralParam.Denom) {
		return sdk.ErrInternal(fmt.Sprintf("param doesnt exists: '%s'", collateralParam.String()))
	}
	for x, cp := range params.CollateralParams {
		if cp.Denom == collateralParam.Denom {
			params.CollateralParams[x] = collateralParam
		}
	}
	k.SetParams(ctx, params)
	return nil
}
