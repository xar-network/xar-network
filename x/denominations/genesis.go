package denominations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

type GenesisState struct {
	TokenRecords []types.Token `json:"token_records"`
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.TokenRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid TokenRecord: Value: %s. Error: Missing Owner", record.Owner)
		}
		if record.Symbol == "" {
			return fmt.Errorf("invalid TokenRecord: Owner: %s. Error: Missing Symbol", record.Symbol)
		}
		if record.TotalSupply == nil || record.TotalSupply.Len() == 0 {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing TotalSupply", record.TotalSupply)
		}
		if record.Name == "" {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing Name", record.Name)
		}
		if record.OriginalSymbol == "" {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing OriginalSymbol", record.OriginalSymbol)
		}
	}
	return nil
}

func NewGenesisState(poa string) GenesisState {
	return GenesisState{TokenRecords: nil}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		TokenRecords: []types.Token{},
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, record := range data.TokenRecords {
		record := record
		err := k.SetToken(ctx, record.Symbol, &record)
		if err != nil {
			panic(fmt.Sprintf("failed to set token for symbol: %s. Error: %s", record.Symbol, err))
		}
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []types.Token
	iterator := k.GetTokensIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		symbol := string(iterator.Key())
		token, err := k.GetToken(ctx, symbol)
		if err != nil {
			panic(fmt.Sprintf("failed to find token for symbol: %s. Error: %s", symbol, err))
		}
		records = append(records, *token)

	}
	return GenesisState{TokenRecords: records}
}
