package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/zar-network/zar-network/x/cdp/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cdpTxCmd := &cobra.Command{
		Use:   "cdp",
		Short: "cdp transactions subcommands",
	}

	cdpTxCmd.AddCommand(
		GetCmdModifyCdp(cdc),
	)

	return cdpTxCmd
}

// GetCmdModifyCdp cli command for creating and modifying cdps.
func GetCmdModifyCdp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modifycdp [from_key_or_addres] [ownerAddress] [collateralType] [collateralChange] [debtChange]",
		Short: "create or modify a cdp",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			collateralChange, ok := sdk.NewIntFromString(args[3])
			if !ok {
				fmt.Printf("invalid collateral amount - %s \n", string(args[3]))
				return nil
			}
			debtChange, ok := sdk.NewIntFromString(args[4])
			if !ok {
				fmt.Printf("invalid debt amount - %s \n", string(args[4]))
				return nil
			}
			msg := types.NewMsgCreateOrModifyCDP(cliCtx.GetFromAddress(), args[2], collateralChange, debtChange)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
