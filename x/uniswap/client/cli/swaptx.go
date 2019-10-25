package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	swap "github.com/zar-network/zar-network/x/uniswap/internal/types"

	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        swap.ModuleName,
		Short:                      "swap subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetSwapTxCmd(cdc),
	)

	return cmd
}

func GetSwapTxCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "swap [asset] [target_denom]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// parse coins
			asset, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()

			msg := swap.NewMsgSwap(from, asset, args[1])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
