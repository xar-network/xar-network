package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	boxcli "github.com/hashgard/hashgard/x/box/client/cli"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/params"
	"github.com/hashgard/hashgard/x/box/types"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCreateCmd implements create Future box transaction command.
func GetCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [name] [total-amount] [distribute-file]",
		Args:    cobra.ExactArgs(3),
		Short:   "Create a new future box",
		Long:    "Create a new future box",
		Example: "$ hashgardcli future create foocoin 100000000coin174876e800 path/distribute.json --from foo",
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
			contents, err := ioutil.ReadFile(args[2])
			if err != nil {
				return err
			}

			futureBox := types.FutureBox{}
			err = json.Unmarshal(contents, &futureBox)
			if err != nil {
				return err
			}
			coin.Amount = boxutils.MulDecimals(coin, decimal)
			if err = processFutureBox(coin, futureBox, decimal); err != nil {
				return err
			}
			box := params.BoxFutureParams{}
			box.Name = args[0]
			box.TotalAmount = types.BoxToken{Token: coin, Decimals: decimal}
			box.TransferDisabled = viper.GetBool(boxcli.FlagTransferDisabled)
			box.Future = futureBox
			msg := msgs.NewMsgFutureBox(account.GetAddress(), &box)
			if err := msg.ValidateService(); err != nil {
				return errors.Errorf(err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	cmd.Flags().Bool(boxcli.FlagTransferDisabled, true, "Disable transfer the box")
	return cmd
}

func processFutureBox(totalAmount sdk.Coin, futureBox types.FutureBox, decimals uint) sdk.Error {
	if futureBox.Receivers == nil {
		return errors.ErrNotSupportOperation()
	}
	total := sdk.ZeroInt()
	for i, items := range futureBox.Receivers {
		for j, rec := range items {
			if j == 0 {
				_, err := sdk.AccAddressFromBech32(rec)
				if err != nil {
					return sdk.ErrInvalidAddress(rec)
				}
				continue
			}
			amount, ok := sdk.NewIntFromString(rec)
			if !ok {
				return errors.ErrAmountNotValid(rec)
			}
			amount = boxutils.MulDecimals(boxutils.ParseCoin(totalAmount.Denom, amount), decimals)
			total = total.Add(amount)
			futureBox.Receivers[i][j] = amount.String()
		}
	}
	if !total.Equal(totalAmount.Amount) {
		return errors.ErrAmountNotValid("Receivers")
	}
	return nil
}
