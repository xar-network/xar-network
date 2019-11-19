package balance

import (
	"github.com/xar-network/xar-network/embedded"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GetQueryRequest struct {
	Address sdk.AccAddress
}

type GetQueryResponseBalance struct {
	Denom  string   `json:"denom"`
	Liquid sdk.Uint `json:"liquid"`
	AtRisk sdk.Uint `json:"at_risk"`
}

type GetQueryResponse struct {
	Balances []GetQueryResponseBalance `json:"balances"`
}

type TransferBalanceRequest struct {
	To     sdk.AccAddress `json:"to"`
	Denom  string         `json:"denom"`
	Amount sdk.Uint       `json:"amount"`
}

type TransferBalanceResponse struct {
	BlockInclusion embedded.BlockInclusion `json:"block_inclusion"`
}
