package compound

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WhoisRecords       []Whois       `json:"whois_records"`
	MoneyMarketRecords []MoneyMarket `json:"moneymarket_records"`
}

func NewGenesisState(MoneyMarketRecords []Whois) GenesisState {
	return GenesisState{MoneyMarketRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.MoneyMarketRecords {
		if record.Owner == nil {
			return fmt.Errorf("Invalid WhoisRecord: Value: %s. Error: Missing Name", record.Name)
		}
		if record.Name == "" {
			return fmt.Errorf("Invalid WhoisRecord: Owner: %s. Error: Missing Owner", record.Owner)
		}
		// if record.InterestRate == '' {
		// 	return fmt.Errorf("Invalid WhoisRecord: Value: %s. Error: Missing InterestRate", record.InterestRate)
		// }
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		MoneyMarketRecords: []MoneyMarket{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.MoneyMarketRecords {
		keeper.SetMarketInfo(ctx, record.Name, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []MoneyMarket
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		var moneymarket MoneyMarket
		moneymarket = k.GetMarketInfo(ctx, name)
		records = append(records, moneymarket)
	}
	return GenesisState{MoneyMarketRecords: records}
}
