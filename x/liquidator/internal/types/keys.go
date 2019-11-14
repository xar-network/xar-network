package types

const (
	// ModuleKey is the name of the module
	ModuleName = "liquidator"
	// StoreKey is the store key string for issue
	StoreKey = ModuleName
	// RouterKey is the message route for issue
	RouterKey = ModuleName
	// QuerierRoute is the querier route for issue
	QuerierRoute = ModuleName
	// Parameter store default namestore
	DefaultParamspace = ModuleName

	QueryGetOutstandingDebt = "outstanding_debt" // Get the outstanding seized debt
)

var ModuleParamsKey = []byte("LiquidatorModuleParams")
