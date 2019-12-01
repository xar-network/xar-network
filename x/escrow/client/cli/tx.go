package cli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/types"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
	"github.com/spf13/cobra"
)

func ProcessBoxDescriptionCmd(cdc *codec.Codec, boxType string, id string, filename string) error {
	if boxutils.GetBoxTypeByValue(id) != boxType {
		return errors.Errorf(errors.ErrUnknownBox(id))
	}
	if err := boxutils.CheckId(id); err != nil {
		return errors.Errorf(err)
	}
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	err = json.Compact(&buffer, contents)
	if err != nil {
		return errors.ErrBoxDescriptionNotValid()
	}
	contents = buffer.Bytes()

	_, err = clientutils.BoxOwnerCheck(cdc, cliCtx, account, id)
	if err != nil {
		return err
	}
	if len(contents) <= 0 || !json.Valid(contents) {
		return errors.ErrBoxDescriptionNotValid()
	}
	msg := msgs.NewMsgBoxDescription(id, account.GetAddress(), contents)

	validateErr := msg.ValidateBasic()

	if validateErr != nil {
		return errors.Errorf(validateErr)
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

func ProcessBoxDisableFeatureCmd(cdc *codec.Codec, boxType string, id string, feature string) error {
	if boxutils.GetBoxTypeByValue(id) != boxType {
		return errors.Errorf(errors.ErrUnknownBox(id))
	}
	_, ok := types.Features[feature]
	if !ok {
		return errors.Errorf(errors.ErrUnknownFeatures())
	}
	if err := boxutils.CheckId(id); err != nil {
		return errors.Errorf(err)
	}
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	boxInfo, err := clientutils.BoxOwnerCheck(cdc, cliCtx, account, id)
	if err != nil {
		return err
	}

	if feature == types.Transfer && boxInfo.GetBoxType() == types.Lock {
		return errors.Errorf(errors.ErrNotSupportOperation())
	}

	msg := msgs.NewMsgBoxDisableFeature(id, account.GetAddress(), feature)
	validateErr := msg.ValidateBasic()
	if validateErr != nil {
		return errors.Errorf(validateErr)
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

// WithdrawCmd implements withdraw a box transaction command.
func WithdrawCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw [box-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Withdraw a box from the account coins",
		Long:    "Box holder withdraw a deposit box or future box from the account coins when the box can be withdraw",
		Example: "$ hashgardcli bank withdraw boxab3jlxpt2ps --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processBoxWithdrawCmd(cdc, args[0])
		},
	}
	cmd = client.PostCommands(cmd)[0]
	_ = cmd.MarkFlagRequired(client.FlagFrom)
	return cmd
}

// ProcessBoxWithdrawCmd implements withdraw a box transaction command.
func processBoxWithdrawCmd(cdc *codec.Codec, id string) error {
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	msg, err := clientutils.GetWithdrawMsg(cdc, cliCtx, account, id)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)

}

func ProcessBoxInject(cdc *codec.Codec, id string, amountStr string, operation string) error {
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	msg, err := clientutils.GetInjectMsg(cdc, cliCtx, account, id, amountStr, operation, true)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			from := cliCtx.GetFromAddress()

			for i, coin := range coins {
				if err = processBoxSend(cdc, cliCtx, &coin); err != nil {
					return err
				}
				if err = processIssueSend(cdc, cliCtx, &coin, from, to); err != nil {
					return err
				}
				coins[i] = coin
			}
			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}
			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.NewMsgSend(from, to, coins, viper.GetString(client.FlagMemo))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	cmd = client.PostCommands(cmd)[0]
	_ = cmd.MarkFlagRequired(client.FlagFrom)
	return cmd
}

func processBoxSend(cdc *codec.Codec, cliCtx context.CLIContext, coin *sdk.Coin) error {
	if !boxutils.IsId(coin.Denom) {
		return nil
	}
	boxInfo, err := boxclientutils.GetBoxByID(cdc, cliCtx, coin.Denom)
	if err != nil {
		return err
	}
	if boxInfo.IsTransferDisabled() {
		return errors.Errorf(errors.ErrCanNotTransfer(coin.Denom))
	}
	if boxInfo.GetBoxType() == types.Future {
		coin.Amount = issueutils.MulDecimals(coin.Amount, boxInfo.GetTotalAmount().Decimals)
	}
	return nil
}
func processIssueSend(cdc *codec.Codec, cliCtx context.CLIContext, coin *sdk.Coin, from sdk.AccAddress, to sdk.AccAddress) error {
	if !issueutils.IsIssueId(coin.Denom) {
		return nil
	}
	issueInfo, err := issueclientutils.GetIssueByID(cdc, cliCtx, coin.Denom)
	if err != nil {
		return err
	}
	coin.Amount = issueutils.MulDecimals(coin.Amount, issueInfo.GetDecimals())
	if err = issueclientutils.CheckFreeze(cdc, cliCtx, issueInfo.GetIssueId(), from, to); err != nil {
		return err
	}
	return nil
}
