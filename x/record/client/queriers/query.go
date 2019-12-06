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

package queriers

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

func GetQueryRecordPath(recordHash string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryRecord, recordHash)
}
func GetQueryParamsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryParams)
}
func GetQueryRecordsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryRecords)
}

func QueryParams(cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryParamsPath(), nil)
}
func QueryRecord(hash string, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryRecordPath(hash), nil)
}
func QueryRecords(params types.RecordQueryParams, cliCtx context.CLIContext) ([]byte, int64, error) {
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, -1, err
	}
	return cliCtx.QueryWithData(GetQueryRecordsPath(), bz)
}
