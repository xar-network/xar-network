//nolint
package app

import (
	"io"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	genesisFile        string
	paramsFile         string
	exportParamsPath   string
	exportParamsHeight int
	exportStatePath    string
	exportStatsPath    string
	seed               int64
	initialBlockHeight int
	numBlocks          int
	blockSize          int
	enabled            bool
	verbose            bool
	lean               bool
	commit             bool
	period             int
	onOperation        bool // TODO Remove in favor of binary search for invariant violation
	allInvariants      bool
	genesisTime        int64
)

// DONTCOVER

// NewZarAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewZarAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (zapp *ZarApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	zapp = NewZarApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return zapp, zapp.keys[baseapp.MainStoreKey], zapp.keys[staking.StoreKey], zapp.stakingKeeper
}
