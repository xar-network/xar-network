package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Liquidator transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmd_SeizeAndStartCollateralAuction(cdc),
		GetCmd_StartDebtAuction(cdc),
	)

	return cmd
}

func GetCmd_SeizeAndStartCollateralAuction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seize [csdt-owner] [collateral-denom]",
		Short: "seize funds from a CSDT and send to auction",
		Long: `Seize a fixed amount of collateral and debt from a CSDT then start an auction with the collateral.
The amount of collateral seized is given by the 'AuctionSize' module parameter or, if there isn't enough collateral in the CSDT, all the CSDT's collateral is seized.
Debt is seized in proportion to the collateral seized so that the CSDT stays at the same collateral to debt ratio.
A 'forward-reverse' auction is started selling the seized collateral for some stable coin, with a maximum bid of stable coin set to equal the debt seized.
As this is a forward-reverse auction type, if the max stable coin is bid then bidding continues by bidding down the amount of collateral taken by the bidder. At the end, extra collateral is returned to the original CSDT owner.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Setup
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			// Validate inputs
			sender := cliCtx.GetFromAddress()
			csdtOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			denom := args[1]
			// TODO validate denom?

			// Prepare and send message
			msgs := []sdk.Msg{types.MsgSeizeAndStartCollateralAuction{
				Sender:          sender,
				CsdtOwner:       csdtOwner,
				CollateralDenom: denom,
			}}
			// TODO print out results like auction ID?
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	return cmd
}

func GetCmd_StartDebtAuction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "start a debt auction, minting gov coin to cover debt",
		Long:  "Start a reverse auction, selling off minted gov coin to raise a fixed amount of stable coin.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Setup
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender := cliCtx.GetFromAddress()

			// Prepare and send message
			msgs := []sdk.Msg{types.MsgStartDebtAuction{
				Sender: sender,
			}}
			// TODO print out results like auction ID?
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	return cmd
}
