package types

const (
	BurnOwner  = "burn-owner"
	BurnHolder = "burn-holder"
	BurnFrom   = "burn-from"
	Freeze     = "freeze"
	Minting    = "minting"
)

var Features = map[string]int{BurnOwner: 1, BurnHolder: 1, BurnFrom: 1, Freeze: 1, Minting: 1}
