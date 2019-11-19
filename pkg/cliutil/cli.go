package cliutil

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func ValidateAndBroadcast(cliCtx context.CLIContext, bldr authtypes.TxBuilder, msg sdk.Msg) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, bldr, []sdk.Msg{msg})
}
