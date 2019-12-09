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
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tdb "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramscutils "github.com/cosmos/cosmos-sdk/x/params/client/utils"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestProposalChangeDepositParams(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	proposal, _ := gapp.paramsKeeper.GetSubspace("gov")

	var depositParams gov.DepositParams

	ctx := gapp.NewContext(true, abci.Header{Height: gapp.LastBlockHeight()})
	proposal.Get(ctx, gov.ParamStoreKeyDepositParams, &depositParams)

	propJson, err := paramscutils.ParseParamChangeProposalJSON(gapp.Codec(), "proposal_deposit_params.json")
	require.NoError(t, err)

	changes := params.NewParameterChangeProposal(propJson.Title, propJson.Description, propJson.Changes.ToParamChanges())
	hdlr := params.NewParamChangeProposalHandler(gapp.paramsKeeper)
	require.NoError(t, hdlr(ctx, changes))

	var newParams gov.DepositParams
	proposal.Get(ctx, gov.ParamStoreKeyDepositParams, &newParams)
	require.Equal(t, newParams.MaxDepositPeriod, time.Duration(172800000000000))
	require.Equal(t, newParams.MinDeposit, sdk.NewCoins(sdk.NewCoin("uftm", sdk.NewInt(1000000000))))
}

func TestProposalChangeInflationRateChange(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	proposal, _ := gapp.paramsKeeper.GetSubspace("mint")

	var inflationRateChange sdk.Dec

	ctx := gapp.NewContext(true, abci.Header{Height: gapp.LastBlockHeight()})
	proposal.Get(ctx, mint.KeyInflationRateChange, &inflationRateChange)
	match, _ := sdk.NewDecFromStr("0.130000000000000000")
	require.Equal(t, inflationRateChange, match)

	propJson, err := paramscutils.ParseParamChangeProposalJSON(gapp.Codec(), "proposal_inflation.json")
	require.NoError(t, err)

	changes := params.NewParameterChangeProposal(propJson.Title, propJson.Description, propJson.Changes.ToParamChanges())
	hdlr := params.NewParamChangeProposalHandler(gapp.paramsKeeper)
	require.NoError(t, hdlr(ctx, changes))

	var newParams sdk.Dec
	proposal.Get(ctx, mint.KeyInflationRateChange, &newParams)
	match, _ = sdk.NewDecFromStr("5.000000000000000000")
	require.Equal(t, newParams, match)
}

func TestProposalChangeVotingPeriod(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	proposal, _ := gapp.paramsKeeper.GetSubspace("gov")

	var votingPeriod gov.VotingParams
	expected := gov.NewVotingParams(172800000000000)

	ctx := gapp.NewContext(true, abci.Header{Height: gapp.LastBlockHeight()})
	proposal.Get(ctx, gov.ParamStoreKeyVotingParams, &votingPeriod)
	require.Equal(t, votingPeriod, expected)

	propJson, err := paramscutils.ParseParamChangeProposalJSON(gapp.Codec(), "proposal_voting_period.json")
	require.NoError(t, err)

	changes := params.NewParameterChangeProposal(propJson.Title, propJson.Description, propJson.Changes.ToParamChanges())
	hdlr := params.NewParamChangeProposalHandler(gapp.paramsKeeper)
	require.NoError(t, hdlr(ctx, changes))

	var newParams gov.VotingParams
	expected = gov.NewVotingParams(604800000000000)
	proposal.Get(ctx, gov.ParamStoreKeyVotingParams, &newParams)
	require.Equal(t, newParams, expected)
}

func TestXardGeneric(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	modAccPerms := GetMaccPerms()
	require.Equal(t, 11, len(modAccPerms))
}

func TestXardValidateGenesis(t *testing.T) {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("xar", "xarp")
	config.SetBech32PrefixForValidator("xva", "xvap")
	config.SetBech32PrefixForConsensusNode("xca", "xcap")
	config.SetKeyringServiceName("xar")
	config.Seal()

	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)
	// Load default if passed no args, otherwise load passed file
	genesis := DefaultNodeHome + "/config/genesises.json"

	t.Logf("validating genesis file at %s\n", genesis)

	genDoc, err := tmtypes.GenesisDocFromFile(genesis)
	if err == nil {
		var genState map[string]json.RawMessage
		if err := gapp.Codec().UnmarshalJSON(genDoc.AppState, &genState); err != nil {
			t.Errorf("error unmarshaling genesis doc %s: %s", genesis, err.Error())
		}

		for _, moduleName := range gapp.MM().OrderInitGenesis {

			err := gapp.MM().Modules[moduleName].ValidateGenesis(genState[moduleName])
			if err != nil {
				if moduleName != "genutil" {
					t.Errorf("error validating genesis file %s[%s]: %s", genesis, moduleName, err.Error())
				}
			}
		}

		// TODO test to make sure initchain doesn't panic

		t.Logf("File at %s is a valid genesis file\n", genesis)
	}
}

func TestXardExport(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestXardExportZeroHeight(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := tdb.NewMemDB()
	mkdb := tdb.NewMemDB()
	gapp := NewXarApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, mkdb, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, gapp.bankKeeper.BlacklistedAddr(gapp.supplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(gapp *XarApp) error {
	//genesisState := simapp.NewDefaultGenesisState()
	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	gapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	gapp.Commit()
	return nil
}
