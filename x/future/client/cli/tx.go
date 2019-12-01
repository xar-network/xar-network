package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	boxcli "github.com/hashgard/hashgard/x/box/client/cli"
	"github.com/hashgard/hashgard/x/box/types"
	"github.com/spf13/cobra"
)

func GetInjectCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inject [id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Deposit token to the future box",
		Long:    "Deposit token to the future box",
		Example: "$ hashgardcli future inject box174876e800 88888 --from foo",
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
		Short:   "Cancel deposit from a future box",
		Long:    "Cancel deposit from a future box",
		Example: "$ hashgardcli future cancel box174876e800 88888 --from foo",
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
		Short:   "Describe a future box",
		Long:    "Box owner can set description of the future box, and the description need to be in json format. You can customize preferences or use recommended templates.",
		Example: "$ hashgardcli future describe boxab3jlxpt2ps path/description.json --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxDescriptionCmd(cdc, types.Future, args[0], args[1])
		},
	}
	return cmd
}

// GetDisableFeatureCmd implements disable feature a box transaction command.
func GetDisableFeatureCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable [id] [feature]",
		Args:  cobra.ExactArgs(2),
		Short: "Disable feature from a future box",
		Long: fmt.Sprintf("Box Owner disabled the features:\n"+
			"%s:Box holder can transfer", types.Transfer),
		Example: fmt.Sprintf("$ hashgardcli future disable boxab3jlxpt2ps %s --from foo", types.Transfer),
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessBoxDisableFeatureCmd(cdc, types.Future, args[0], args[1])
		},
	}
	return cmd
}
