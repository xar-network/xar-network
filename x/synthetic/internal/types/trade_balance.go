package types

type TotalTradeBalance struct {
	Markets []*MarketBalance `json:"markets" yaml:"markets"`
}

func (t TotalTradeBalance) SetMarketBalance(denom string, newBalance *MarketBalance) {
	mbId, found := t.getMarketBalanceId(denom)
	if found {
		t.Markets[mbId] = newBalance
		return
	}

	t.Markets = append(t.Markets, newBalance)
}

func (t TotalTradeBalance) GetMarketBalance(denom string) *MarketBalance {
	for _, v := range t.Markets {
		if v.MarketDenom == denom {
			return v
		}
	}
	return nil
}

func (t TotalTradeBalance) getMarketBalanceId(denom string) (int, bool) {
	for i, v := range t.Markets {
		if v.MarketDenom == denom {
			return i, true
		}
	}
	return 0, false
}
