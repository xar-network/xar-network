package types

type MarketBalances []MarketBalance

func (m MarketBalances) GetMarketBalance(denom string) (MarketBalance, bool) {
	for _, v := range m {
		if v.MarketDenom == denom {
			return v, true
		}
	}
	return MarketBalance{}, false
}

func (m *MarketBalances) SetMarketBalance(mb MarketBalance) {
	for i, v := range *m {
		if v.MarketDenom == mb.MarketDenom {
			(*m)[i] = mb
			return
		}
	}
	*m = append(*m, mb)
}