package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

const (
	QueryParams    = "params"
	QueryIssues    = "list"
	QueryIssue     = "query"
	QueryAllowance = "allowance"
	QueryFreeze    = "freeze"
	QueryFreezes   = "freezes"
	QuerySearch    = "search"
)

//NewQuerier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, k)
		case QueryIssue:
			return queryIssue(ctx, path[1], k)
		case QueryAllowance:
			return queryAllowance(ctx, path[1], path[2], path[3], k)
		case QueryFreeze:
			return queryFreeze(ctx, path[1], path[2], k)
		case QueryFreezes:
			return queryFreezes(ctx, path[1], k)
		case QuerySearch:
			return querySymbol(ctx, path[1], k)
		case QueryIssues:
			return queryIssues(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown issue query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryIssue(ctx sdk.Context, issueID string, k Keeper) ([]byte, sdk.Error) {
	issue := k.GetIssue(ctx, issueID)
	if issue == nil {
		return nil, types.ErrUnknownIssue()
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, issue)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryAllowance(
	ctx sdk.Context,
	issueID string,
	fromAddr string,
	toAddr string,
	k Keeper,
) ([]byte, sdk.Error) {
	fromAddress, _ := sdk.AccAddressFromBech32(fromAddr)
	toAddress, _ := sdk.AccAddressFromBech32(toAddr)
	amount := k.Allowance(ctx, fromAddress, toAddress, issueID)

	if amount.GT(sdk.ZeroInt()) {
		coinIssueInfo := k.GetIssue(ctx, issueID)
		amount = types.QuoDecimals(amount, coinIssueInfo.GetDecimals())
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, types.NewApproval(amount))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryFreeze(ctx sdk.Context, issueID string, accAddress string, k Keeper) ([]byte, sdk.Error) {
	address, _ := sdk.AccAddressFromBech32(accAddress)
	freeze := k.GetFreeze(ctx, address, issueID)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, freeze)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryFreezes(ctx sdk.Context, issueID string, k Keeper) ([]byte, sdk.Error) {
	freeze := k.GetFreezes(ctx, issueID)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, freeze)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySymbol(ctx sdk.Context, symbol string, k Keeper) ([]byte, sdk.Error) {
	issue := k.SearchIssues(ctx, symbol)
	if issue == nil {
		return nil, types.ErrUnknownIssue()
	}
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, issue)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryIssues(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.IssueQueryParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	issues := k.List(ctx, params)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, issues)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
