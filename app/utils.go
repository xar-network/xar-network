//nolint
package app

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// ExportStateToJSON util function to export the app state to JSON
func ExportStateToJSON(app *xarApp, path string) error {
	fmt.Println("exporting app state...")
	appState, _, err := app.ExportAppStateAndValidators(false, nil)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(appState), 0644)
}

// NewxarAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewxarAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (app *xarApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	app = NewXarApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return app, app.keys[bam.MainStoreKey], app.keys[staking.StoreKey], app.stakingKeeper
}
