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
