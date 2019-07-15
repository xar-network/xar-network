//nolint
package app

import (
	"io"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/Fantom-foundation/cosmos-sdk/baseapp"
	sdk "github.com/Fantom-foundation/cosmos-sdk/types"
	"github.com/Fantom-foundation/cosmos-sdk/x/staking"
)

// DONTCOVER

// NewZarAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewZarAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) (gapp *ZarApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	gapp = NewZarApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return gapp, gapp.keyMain, gapp.keyStaking, gapp.stakingKeeper
}
