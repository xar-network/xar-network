package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	QueryGetCsdts             = "cdts"
	QueryGetParams            = "params"
	RestOwner                 = "owner"
	RestCollateralDenom       = "collateralDenom"
	RestUnderCollateralizedAt = "underCollateralizedAt"
)

type QueryCsdtsParams struct {
	CollateralDenom       string         // get CSDTs with this collateral denom
	Owner                 sdk.AccAddress // get CSDTs belonging to this owner
	UnderCollateralizedAt sdk.Dec        // get CSDTs that will be below the liquidation ratio when the collateral is at this price.
}

type ModifyCsdtRequestBody struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Csdt    CSDT         `json:"csdt"`
}
