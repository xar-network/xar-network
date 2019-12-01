package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xar-network/xar-network/x/escrow/internal/keeper"
	"github.com/xar-network/xar-network/x/escrow/internal/types"
)

const (
	FlagAddress          = "address"
	FlagLimit            = "limit"
	FlagStartId          = "start-id"
	FlagTransferDisabled = "transfer-disabled"

	FlagBottomLine    = "bottom-line"
	FlagPrice         = "price"
	FlagInterest      = "interest"
	FlagStartTime     = "start-time"
	FlagEstablishTime = "establish-time"
	FlagMaturityTime  = "maturity-time"
)

//Box query
// QueryCmd implements the query box command.
func QueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "box [box-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query the details of the account box",
		Long:    "Query the details of the account box",
		Example: "$ xarcli bank box boxab3jlxpt2ps",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processQueryBoxCmd(cdc, args[0])
		},
	}
}

// ProcessQueryBoxParamsCmd implements the query box params command.
func ProcessQueryBoxParamsCmd(cdc *codec.Codec, boxType string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	res, err := keeper.QueryBoxParams(cliCtx)
	if err != nil {
		return err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(res, &params)
	return cliCtx.PrintOutput(types.GetBoxParams(params, boxType))
}

// ProcessQueryBoxCmd implements the query box command.
func ProcessQueryBoxCmd(cdc *codec.Codec, boxType string, id string) error {
	if boxtypes.GetBoxTypeByValue(id) != boxType {
		return types.Errorf(types.ErrUnknownBox(id))
	}
	return processQueryBoxCmd(cdc, id)
}

// ProcessQueryBoxCmd implements the query box command.
func processQueryBoxCmd(cdc *codec.Codec, id string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	if err := boxtypes.CheckId(id); err != nil {
		return types.Errorf(err)
	}
	// Query the box
	res, err := boxkeeper.QueryBoxByID(id, cliCtx)
	if err != nil {
		return err
	}
	var box types.BoxInfo
	cdc.MustUnmarshalJSON(res, &box)
	return cliCtx.PrintOutput(types.GetBoxInfo(box))
}

// ProcessListBoxCmd implements the query box command.
func ProcessListBoxCmd(cdc *codec.Codec, boxType string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	_, ok := types.BoxType[boxType]
	if !ok {
		return types.Errorf(types.ErrUnknownBoxType())
	}
	address, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddress))
	if err != nil {
		return err
	}
	boxQueryParams := types.BoxQueryParams{
		StartId: viper.GetString(FlagStartId),
		BoxType: boxType,
		Owner:   address,
		Limit:   viper.GetInt(FlagLimit),
	}
	// Query the box
	res, err := boxkeeper.QueryBoxsList(boxQueryParams, cdc, cliCtx)
	if err != nil {
		return err
	}
	var boxs types.BoxInfos
	cdc.MustUnmarshalJSON(res, &boxs)
	return cliCtx.PrintOutput(types.GetBoxList(boxs, boxQueryParams.BoxType))
}

// ProcessSearchBoxsCmd implements the query box command.
func ProcessSearchBoxsCmd(cdc *codec.Codec, boxType string, name string) error {
	_, ok := types.BoxType[boxType]

	if !ok {
		return types.Errorf(types.ErrUnknownBoxType())
	}
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	// Query the box
	res, err := boxkeeper.QueryBoxByName(boxType, strings.ToLower(name), cliCtx)
	if err != nil {
		return err
	}
	var boxs types.BoxInfos
	cdc.MustUnmarshalJSON(res, &boxs)
	return cliCtx.PrintOutput(types.GetBoxList(boxs, boxType))
}
