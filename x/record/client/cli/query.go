package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	recordqueriers "github.com/xar-network/xar-network/x/record/client/queriers"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

// GetCmdQueryRecord implements the query record command.
func GetCmdQueryRecord(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query [hash]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a single record",
		Long:    "Query detail of a record by record hash",
		Example: "$ xarcli record query BC38CAEE32149BEF4CCFAEAB518EC9A5FBC85AE6AC8D5A9F6CD710FAF5E4A2B8",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processQuery(cdc, args)
		},
	}
}

func processQuery(cdc *codec.Codec, args []string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	hash := args[0]
	if err := types.CheckRecordHash(hash); err != nil {
		return types.Errorf(err)
	}
	// Query the record
	res, height, err := recordqueriers.QueryRecord(hash, cliCtx)
	if err != nil {
		return err
	}
	cliCtx = cliCtx.WithHeight(height)

	var recordInfo types.Record
	cdc.MustUnmarshalJSON(res, &recordInfo)
	return cliCtx.PrintOutput(recordInfo)
	//_, err = cliCtx.Output.Write(res)
	//if err != nil {
	//	return err
	//}
	//return nil
}

// GetCmdQueryList implements the query records command.
func GetCmdQueryList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Query record list",
		Long:    "Query record list, flag sender is optional, default limit is 30",
		Example: "$ xarcli record list --sender xar1s6auwlcevspesynw44vx23e3zhuz7as9ulz56l",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender, err := sdk.AccAddressFromBech32(viper.GetString(flagSender))
			if err != nil {
				return err
			}
			startId := viper.GetString(flagStartRecordId)
			if len(startId) > 0 {
				if err := types.CheckRecordId(startId); err != nil {
					return types.Errorf(err)
				}
			}
			recordQueryParams := types.RecordQueryParams{
				StartRecordId: startId,
				Sender:        sender,
				Limit:         30,
			}
			limit := viper.GetInt(flagLimit)
			if limit > 0 {
				recordQueryParams.Limit = limit
			}
			// Query the record
			res, height, err := recordqueriers.QueryRecords(recordQueryParams, cliCtx)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			var ls types.Records
			cdc.MustUnmarshalJSON(res, &ls)
			if len(ls) == 0 {
				fmt.Println("No records")
				return nil
			}
			return cliCtx.PrintOutput(ls)
			//_, err = cliCtx.Output.Write(res)
			//if err != nil {
			//	return err
			//}
			//return nil
		},
	}

	cmd.Flags().String(flagSender, "", "sender address")
	cmd.Flags().String(flagStartRecordId, "", "Start recordId of records")
	cmd.Flags().Int32(flagLimit, 30, "Query number of record results per page returned")
	return cmd
}
