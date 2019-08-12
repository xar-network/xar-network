package queriers

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

func GetQueryIssuePath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryIssue, issueID)
}
func GetQueryParamsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryParams)
}
func GetQueryIssueAllowancePath(issueID string, owner sdk.AccAddress, spender sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryAllowance, issueID, owner.String(), spender.String())
}
func GetQueryIssueFreezePath(issueID string, accAddress sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryFreeze, issueID, accAddress.String())
}
func GetQueryIssueFreezesPath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryFreezes, issueID)
}
func GetQueryIssueSearchPath(symbol string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QuerySearch, symbol)
}
func GetQueryIssuesPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryIssues)
}

func QueryIssueBySymbol(symbol string, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryIssueSearchPath(symbol), nil)
}
func QueryParams(cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryParamsPath(), nil)
}
func QueryIssueByID(issueID string, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryIssuePath(issueID), nil)
}
func QueryIssueAllowance(issueID string, owner sdk.AccAddress, spender sdk.AccAddress, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryIssueAllowancePath(issueID, owner, spender), nil)
}
func QueryIssueFreeze(issueID string, accAddress sdk.AccAddress, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryIssueFreezePath(issueID, accAddress), nil)
}
func QueryIssueFreezes(issueID string, cliCtx context.CLIContext) ([]byte, int64, error) {
	return cliCtx.QueryWithData(GetQueryIssueFreezesPath(issueID), nil)
}
func QueryIssuesList(params types.IssueQueryParams, cdc *codec.Codec, cliCtx context.CLIContext) ([]byte, int64, error) {
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, -1, err
	}
	return cliCtx.QueryWithData(GetQueryIssuesPath(), bz)
}
