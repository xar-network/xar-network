package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	boxcli "github.com/hashgard/hashgard/x/box/client/cli"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/types"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
	"github.com/spf13/cobra"
)

// GetInterestInjectCmd implements interest injection a deposit box transaction command.
func GetInterestInjectCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interest-inject [id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Inject interest to the deposit box",
		Long:    "Inject interest to the deposit box",
		Example: "$ hashgardcli deposit interest-inject box174876e800 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return interest(cdc, args[0], args[1], types.Inject)
		},
	}
	return cmd
}

// GetInterestCancelCmd implements fetch interest from a deposit box transaction command.
func GetInterestCancelCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interest-cancel [id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Cancel interest from a deposit box",
		Long:    "Cancel interest from a deposit box",
		Example: "$ hashgardcli deposit interest-cancel box174876e800 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return interest(cdc, args[0], args[1], types.Cancel)
		},
	}
	return cmd
}

func interest(cdc *codec.Codec, id string, amountStr string, operation string) error {
	if boxutils.GetBoxTypeByValue(id) != types.Deposit {
		return errors.Errorf(errors.ErrNotSupportOperation())
	}
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	msg, err := clientutils.GetInterestMsg(cdc, cliCtx, account, id, amountStr, operation, true)
	if err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

func GetInjectCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inject [id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Deposit to the deposit box",
		Long:    "Deposit to the deposit box",
		Example: "$ hashgardcli deposit deposit-to box174876e800 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxInject(cdc, args[0], args[1], types.Inject)
		},
	}
	return cmd
}

func GetCancelDepositCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cancel [id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Cancel deposit from a deposit box",
		Long:    "Cancel deposit from a deposit box",
		Example: "$ hashgardcli deposit fetch box174876e800 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxInject(cdc, args[0], args[1], types.Cancel)
		},
	}
	return cmd
}

func GetDescriptionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe [id] [description-file]",
		Args:    cobra.ExactArgs(2),
		Short:   "Describe a deposit box",
		Long:    "Box owner can set description of the deposit box, and the description need to be in json format. You can customize preferences or use recommended templates.",
		Example: "$ hashgardcli deposit describe boxab3jlxpt2ps path/description.json --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxDescriptionCmd(cdc, types.Deposit, args[0], args[1])
		},
	}
	return cmd
}

// GetDisableFeatureCmd implements disable feature a box transaction command.
func GetDisableFeatureCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable [id] [feature]",
		Args:  cobra.ExactArgs(2),
		Short: "Disable feature from a deposit box",
		Long: fmt.Sprintf("Box Owner disabled the features:\n"+
			"%s:Box holder can transfer", types.Transfer),
		Example: fmt.Sprintf("$ hashgardcli deposit disable boxab3jlxpt2ps %s --from foo", types.Transfer),
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxDisableFeatureCmd(cdc, types.Deposit, args[0], args[1])
		},
	}
	return cmd
}
