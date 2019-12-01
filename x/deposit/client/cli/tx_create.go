package cli

import (
	"github.com/hashgard/hashgard/x/box/params"

	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	boxcli "github.com/hashgard/hashgard/x/box/client/cli"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/types"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCreateCmd implements create deposit box transaction command.
func GetCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [name] [total-amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Create a new deposit box",
		Long:    "Create a new deposit box",
		Example: "$ hashgardcli deposit create foocoin 100000000coin174876e800 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse coins trying to be sent
			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
			if err != nil {
				return err
			}

			decimal, err := clientutils.GetCoinDecimal(cdc, cliCtx, coin)
			if err != nil {
				return err
			}
			coin.Amount = boxutils.MulDecimals(coin, decimal)

			box := params.BoxDepositParams{}
			box.Name = args[0]
			box.TotalAmount = types.BoxToken{Token: coin, Decimals: decimal}
			box.TransferDisabled = viper.GetBool(boxcli.FlagTransferDisabled)
			box.Deposit = types.DepositBox{
				StartTime:     viper.GetInt64(boxcli.FlagStartTime),
				EstablishTime: viper.GetInt64(boxcli.FlagEstablishTime),
				MaturityTime:  viper.GetInt64(boxcli.FlagMaturityTime)}

			num, ok := sdk.NewIntFromString(viper.GetString(boxcli.FlagBottomLine))
			if !ok {
				return errors.Errorf(errors.ErrAmountNotValid(boxcli.FlagBottomLine))
			}
			box.Deposit.BottomLine = num
			num, ok = sdk.NewIntFromString(viper.GetString(boxcli.FlagPrice))
			if !ok {
				return errors.Errorf(errors.ErrAmountNotValid(boxcli.FlagPrice))
			}
			box.Deposit.Price = num
			box.Deposit.Price = boxutils.MulDecimals(boxutils.ParseCoin(box.TotalAmount.Token.Denom, box.Deposit.Price), decimal)
			box.Deposit.BottomLine = boxutils.MulDecimals(boxutils.ParseCoin(box.TotalAmount.Token.Denom, box.Deposit.BottomLine), decimal)

			interest, err := sdk.ParseCoin(viper.GetString(boxcli.FlagInterest))
			if err != nil {
				return err
			}
			decimal, err = clientutils.GetCoinDecimal(cdc, cliCtx, interest)
			if err != nil {
				return err
			}

			interest.Amount = boxutils.MulDecimals(interest, decimal)
			box.Deposit.Interest = types.BoxToken{Token: interest, Decimals: decimal}

			box.Deposit.PerCoupon = boxutils.CalcInterestRate(box.TotalAmount.Token.Amount, box.Deposit.Price,
				box.Deposit.Interest.Token, box.Deposit.Interest.Decimals)

			msg := msgs.NewMsgDepositBox(account.GetAddress(), &box)
			if err := msg.ValidateService(); err != nil {
				return errors.Errorf(err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	cmd.Flags().Bool(boxcli.FlagTransferDisabled, true, "Disable the box transfer")
	cmd.Flags().String(boxcli.FlagBottomLine, "", "Box bottom line")
	cmd.Flags().String(boxcli.FlagPrice, "", "Box unit price")
	cmd.Flags().String(boxcli.FlagInterest, "", "Box interest")
	cmd.Flags().Int64(boxcli.FlagStartTime, 0, "Box start time")
	cmd.Flags().Int64(boxcli.FlagEstablishTime, 0, "Box establish time")
	cmd.Flags().Int64(boxcli.FlagMaturityTime, 0, "Box maturity time")

	return cmd
}

func MarkCmdDepositBoxCreateFlagRequired(cmd *cobra.Command) {
	cmd.MarkFlagRequired(boxcli.FlagBottomLine)
	cmd.MarkFlagRequired(boxcli.FlagPrice)
	cmd.MarkFlagRequired(boxcli.FlagInterest)
	cmd.MarkFlagRequired(boxcli.FlagStartTime)
	cmd.MarkFlagRequired(boxcli.FlagEstablishTime)
	cmd.MarkFlagRequired(boxcli.FlagMaturityTime)
}
