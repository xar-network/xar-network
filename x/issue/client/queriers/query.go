package queriers

/*func QueryIssueBySymbol(symbol string, cliCtx context.CLIContext) ([]byte, int64, error) {
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
}*/
