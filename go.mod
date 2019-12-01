module github.com/xar-network/xar-network

go 1.13

require (
	github.com/99designs/keyring v1.1.3
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d
	github.com/cosmos/cosmos-sdk v0.34.4-0.20191114141721-d4c831e63ad3
	github.com/cosmos/modules/incubator/nft v0.0.0-20191031141754-76e462805729
	github.com/go-kit/kit v0.9.0
	github.com/gobuffalo/packr v1.30.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/sessions v1.1.3
	github.com/olekukonko/tablewriter v0.0.2
	github.com/otiai10/copy v1.0.2
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/common v0.6.0
	github.com/rs/cors v1.7.0
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
)

//replace github.com/cosmos/cosmos-sdk => github.com/Fantom-foundation/cosmos-sdk v0.28.2-0.20190715083311-79f8f2370f7e
