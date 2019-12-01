package cli

import (
	"strconv"

	"github.com/hashgard/hashgard/x/box/params"

	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
	"github.com/spf13/cobra"

	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/types"
)

// GetCreate implements create lock transaction command.
func GetCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [name] [total-amount] [end-time]",
		Args:    cobra.ExactArgs(3),
		Short:   "Create a new lock",
		Long:    "Create a new lock",
		Example: "$ hashgardcli lock create foocoin 100000000coin174876e800 2557223200 --from foo",
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

			endTime, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			coin.Amount = boxutils.MulDecimals(coin, decimal)
			box := &params.BoxLockParams{}
			box.Name = args[0]
			box.TotalAmount = types.BoxToken{Token: coin, Decimals: decimal}
			box.Lock = types.LockBox{EndTime: endTime}

			msg := msgs.NewMsgLockBox(account.GetAddress(), box)
			if err := msg.ValidateService(); err != nil {
				return errors.Errorf(err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	return cmd
}
