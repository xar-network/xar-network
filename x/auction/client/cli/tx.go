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
	"github.com/zar-network/zar-network/x/auction/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auction transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdPlaceBid(cdc),
	)

	return cmd
}

// GetCmdPlaceBid cli command for creating and modifying cdps.
func GetCmdPlaceBid(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "placebid [from_key_or_address] [AuctionID] [Bidder] [Bid] [Lot]",
		Short: "place a bid on an auction",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			id, err := types.NewIDFromString(args[1])
			if err != nil {
				fmt.Printf("invalid auction id - %s \n", string(args[1]))
				return err
			}

			bid, err := sdk.ParseCoin(args[3])
			if err != nil {
				fmt.Printf("invalid bid amount - %s \n", string(args[3]))
				return err
			}

			lot, err := sdk.ParseCoin(args[4])
			if err != nil {
				fmt.Printf("invalid lot - %s \n", string(args[4]))
				return err
			}
			msg := types.NewMsgPlaceBid(id, cliCtx.GetFromAddress(), bid, lot)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
