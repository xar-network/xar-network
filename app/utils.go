//nolint
/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

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
func ExportStateToJSON(app *XarApp, path string) error {
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
func NewXarAppUNSAFE(logger log.Logger, db dbm.DB, mkdb dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (app *XarApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {

	app = NewXarApp(logger, db, mkdb, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return app, app.keys[bam.MainStoreKey], app.keys[staking.StoreKey], app.stakingKeeper
}
