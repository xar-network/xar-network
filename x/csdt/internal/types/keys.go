package types

const (
	// ModuleKey is the name of the module
	ModuleName = "csdt"
	// StoreKey is the store key string for issue
	StoreKey = ModuleName
	// RouterKey is the message route for issue
	RouterKey = ModuleName
	// QuerierRoute is the querier route for issue
	QuerierRoute = ModuleName
	// Parameter store default namestore
	DefaultParamspace = ModuleName

	QueryGetCsdts  = "csdts"
	QueryGetParams = "params"

	// StableDenom asset code of the dollar-denominated debt coin
	StableDenom = "csdt" // TODO allow to be changed
	// GovDenom asset code of the governance coin
	GovDenom = "ftm"
)
