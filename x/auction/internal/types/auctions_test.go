/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

const DefaultMaxBidDuration types.EndTime = 3 * 1

// TODO can this be less verbose? Should PlaceBid() be split into smaller functions?
// It would be possible to combine all auction tests into one test runner.
func TestForwardAuction_PlaceBid(t *testing.T) {
	seller := sdk.AccAddress([]byte("a_seller"))
	buyer1 := sdk.AccAddress([]byte("buyer1"))
	buyer2 := sdk.AccAddress([]byte("buyer2"))
	end := types.EndTime(13)
	now := types.EndTime(10)

	type args struct {
		currentBlockHeight types.EndTime
		bidder             sdk.AccAddress
		lot                sdk.Coin
		bid                sdk.Coin
	}
	tests := []struct {
		name            string
		auction         types.ForwardAuction
		args            args
		expectedOutputs []types.BankOutput
		expectedInputs  []types.BankInput
		expectedEndTime types.EndTime
		expectedBidder  sdk.AccAddress
		expectedBid     sdk.Coin
		expectpass      bool
	}{
		{
			"normal",
			types.ForwardAuction{types.BaseAuction{
				Initiator:  seller,
				Lot:        c("usdx", 100),
				Bidder:     buyer1,
				Bid:        c("ftm", 6),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, buyer2, c("usdx", 100), c("ftm", 10)},
			[]types.BankOutput{{buyer2, c("ftm", 10)}},
			[]types.BankInput{{buyer1, c("ftm", 6)}, {seller, c("ftm", 4)}},
			now + DefaultMaxBidDuration,
			buyer2,
			c("ftm", 10),
			true,
		},
		{
			"lowBid",
			types.ForwardAuction{types.BaseAuction{
				Initiator:  seller,
				Lot:        c("usdx", 100),
				Bidder:     buyer1,
				Bid:        c("ftm", 6),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, buyer2, c("usdx", 100), c("ftm", 5)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			buyer1,
			c("ftm", 6),
			false,
		},
		{
			"equalBid",
			types.ForwardAuction{types.BaseAuction{
				Initiator:  seller,
				Lot:        c("usdx", 100),
				Bidder:     buyer1,
				Bid:        c("ftm", 6),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, buyer2, c("usdx", 100), c("ftm", 6)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			buyer1,
			c("ftm", 6),
			false,
		},
		{
			"timeout",
			types.ForwardAuction{types.BaseAuction{
				Initiator:  seller,
				Lot:        c("usdx", 100),
				Bidder:     buyer1,
				Bid:        c("ftm", 6),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{end + 1, buyer2, c("usdx", 100), c("ftm", 10)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			buyer1,
			c("ftm", 6),
			false,
		},
		{
			"hitMaxEndTime",
			types.ForwardAuction{types.BaseAuction{
				Initiator:  seller,
				Lot:        c("usdx", 100),
				Bidder:     buyer1,
				Bid:        c("ftm", 6),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{end - 1, buyer2, c("usdx", 100), c("ftm", 10)},
			[]types.BankOutput{{buyer2, c("ftm", 10)}},
			[]types.BankInput{{buyer1, c("ftm", 6)}, {seller, c("ftm", 4)}},
			end, // end time should be capped at MaxEndTime
			buyer2,
			c("ftm", 10),
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// update auction and return in/outputs
			outputs, inputs, err := tc.auction.PlaceBid(tc.args.currentBlockHeight, tc.args.bidder, tc.args.lot, tc.args.bid)

			// check for err
			if tc.expectpass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
			// check for correct in/outputs
			require.Equal(t, tc.expectedOutputs, outputs)
			require.Equal(t, tc.expectedInputs, inputs)
			// check for correct EndTime, bidder, bid
			require.Equal(t, tc.expectedEndTime, tc.auction.EndTime)
			require.Equal(t, tc.expectedBidder, tc.auction.Bidder)
			require.Equal(t, tc.expectedBid, tc.auction.Bid)
		})
	}
}

func TestReverseAuction_PlaceBid(t *testing.T) {
	buyer := sdk.AccAddress([]byte("a_buyer"))
	seller1 := sdk.AccAddress([]byte("seller1"))
	seller2 := sdk.AccAddress([]byte("seller2"))
	end := types.EndTime(13)
	now := types.EndTime(10)

	type args struct {
		currentBlockHeight types.EndTime
		bidder             sdk.AccAddress
		lot                sdk.Coin
		bid                sdk.Coin
	}
	tests := []struct {
		name            string
		auction         types.ReverseAuction
		args            args
		expectedOutputs []types.BankOutput
		expectedInputs  []types.BankInput
		expectedEndTime types.EndTime
		expectedBidder  sdk.AccAddress
		expectedLot     sdk.Coin
		expectpass      bool
	}{
		{
			"normal",
			types.ReverseAuction{types.BaseAuction{
				Initiator:  buyer,
				Lot:        c("ftm", 10),
				Bidder:     seller1,
				Bid:        c("usdx", 100),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, seller2, c("ftm", 9), c("usdx", 100)},
			[]types.BankOutput{{seller2, c("usdx", 100)}},
			[]types.BankInput{{seller1, c("usdx", 100)}, {buyer, c("ftm", 1)}},
			now + DefaultMaxBidDuration,
			seller2,
			c("ftm", 9),
			true,
		},
		{
			"highBid",
			types.ReverseAuction{types.BaseAuction{
				Initiator:  buyer,
				Lot:        c("ftm", 10),
				Bidder:     seller1,
				Bid:        c("usdx", 100),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, seller2, c("ftm", 11), c("usdx", 100)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			seller1,
			c("ftm", 10),
			false,
		},
		{
			"equalBid",
			types.ReverseAuction{types.BaseAuction{
				Initiator:  buyer,
				Lot:        c("ftm", 10),
				Bidder:     seller1,
				Bid:        c("usdx", 100),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{now, seller2, c("ftm", 10), c("usdx", 100)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			seller1,
			c("ftm", 10),
			false,
		},
		{
			"timeout",
			types.ReverseAuction{types.BaseAuction{
				Initiator:  buyer,
				Lot:        c("ftm", 10),
				Bidder:     seller1,
				Bid:        c("usdx", 100),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{end + 1, seller2, c("ftm", 9), c("usdx", 100)},
			[]types.BankOutput{},
			[]types.BankInput{},
			end,
			seller1,
			c("ftm", 10),
			false,
		},
		{
			"hitMaxEndTime",
			types.ReverseAuction{types.BaseAuction{
				Initiator:  buyer,
				Lot:        c("ftm", 10),
				Bidder:     seller1,
				Bid:        c("usdx", 100),
				EndTime:    end,
				MaxEndTime: end,
			}},
			args{end - 1, seller2, c("ftm", 9), c("usdx", 100)},
			[]types.BankOutput{{seller2, c("usdx", 100)}},
			[]types.BankInput{{seller1, c("usdx", 100)}, {buyer, c("ftm", 1)}},
			end, // end time should be capped at MaxEndTime
			seller2,
			c("ftm", 9),
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// update auction and return in/outputs
			outputs, inputs, err := tc.auction.PlaceBid(tc.args.currentBlockHeight, tc.args.bidder, tc.args.lot, tc.args.bid)

			// check for err
			if tc.expectpass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
			// check for correct in/outputs
			require.Equal(t, tc.expectedOutputs, outputs)
			require.Equal(t, tc.expectedInputs, inputs)
			// check for correct EndTime, bidder, bid
			require.Equal(t, tc.expectedEndTime, tc.auction.EndTime)
			require.Equal(t, tc.expectedBidder, tc.auction.Bidder)
			require.Equal(t, tc.expectedLot, tc.auction.Lot)
		})
	}
}

func TestForwardReverseAuction_PlaceBid(t *testing.T) {
	cdpOwner := sdk.AccAddress([]byte("a_cdp_owner"))
	seller := sdk.AccAddress([]byte("a_seller"))
	buyer1 := sdk.AccAddress([]byte("buyer1"))
	buyer2 := sdk.AccAddress([]byte("buyer2"))
	end := types.EndTime(13)
	now := types.EndTime(10)

	type args struct {
		currentBlockHeight types.EndTime
		bidder             sdk.AccAddress
		lot                sdk.Coin
		bid                sdk.Coin
	}
	tests := []struct {
		name            string
		auction         types.ForwardReverseAuction
		args            args
		expectedOutputs []types.BankOutput
		expectedInputs  []types.BankInput
		expectedEndTime types.EndTime
		expectedBidder  sdk.AccAddress
		expectedLot     sdk.Coin
		expectedBid     sdk.Coin
		expectpass      bool
	}{
		{
			"normalForwardBid",
			types.ForwardReverseAuction{BaseAuction: types.BaseAuction{
				Initiator:  seller,
				Lot:        c("xrp", 100),
				Bidder:     buyer1,
				Bid:        c("usdx", 5),
				EndTime:    end,
				MaxEndTime: end},
				MaxBid:      c("usdx", 10),
				OtherPerson: cdpOwner,
			},
			args{now, buyer2, c("xrp", 100), c("usdx", 6)},
			[]types.BankOutput{{buyer2, c("usdx", 6)}},
			[]types.BankInput{{buyer1, c("usdx", 5)}, {seller, c("usdx", 1)}},
			now + DefaultMaxBidDuration,
			buyer2,
			c("xrp", 100),
			c("usdx", 6),
			true,
		},
		{
			"normalSwitchOverBid",
			types.ForwardReverseAuction{BaseAuction: types.BaseAuction{
				Initiator:  seller,
				Lot:        c("xrp", 100),
				Bidder:     buyer1,
				Bid:        c("usdx", 5),
				EndTime:    end,
				MaxEndTime: end},
				MaxBid:      c("usdx", 10),
				OtherPerson: cdpOwner,
			},
			args{now, buyer2, c("xrp", 99), c("usdx", 10)},
			[]types.BankOutput{{buyer2, c("usdx", 10)}},
			[]types.BankInput{{buyer1, c("usdx", 5)}, {seller, c("usdx", 5)}, {cdpOwner, c("xrp", 1)}},
			now + DefaultMaxBidDuration,
			buyer2,
			c("xrp", 99),
			c("usdx", 10),
			true,
		},
		{
			"normalReverseBid",
			types.ForwardReverseAuction{BaseAuction: types.BaseAuction{
				Initiator:  seller,
				Lot:        c("xrp", 99),
				Bidder:     buyer1,
				Bid:        c("usdx", 10),
				EndTime:    end,
				MaxEndTime: end},
				MaxBid:      c("usdx", 10),
				OtherPerson: cdpOwner,
			},
			args{now, buyer2, c("xrp", 90), c("usdx", 10)},
			[]types.BankOutput{{buyer2, c("usdx", 10)}},
			[]types.BankInput{{buyer1, c("usdx", 10)}, {cdpOwner, c("xrp", 9)}},
			now + DefaultMaxBidDuration,
			buyer2,
			c("xrp", 90),
			c("usdx", 10),
			true,
		},
		// TODO more test cases
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// update auction and return in/outputs
			outputs, inputs, err := tc.auction.PlaceBid(tc.args.currentBlockHeight, tc.args.bidder, tc.args.lot, tc.args.bid)

			// check for err
			if tc.expectpass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
			// check for correct in/outputs
			require.Equal(t, tc.expectedOutputs, outputs)
			require.Equal(t, tc.expectedInputs, inputs)
			// check for correct EndTime, bidder, bid
			require.Equal(t, tc.expectedEndTime, tc.auction.EndTime)
			require.Equal(t, tc.expectedBidder, tc.auction.Bidder)
			require.Equal(t, tc.expectedLot, tc.auction.Lot)
			require.Equal(t, tc.expectedBid, tc.auction.Bid)
		})
	}
}

// defined to avoid cluttering test cases with long function name
func c(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}
