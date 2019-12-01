package tests

import (
	"testing"

	"github.com/hashgard/hashgard/x/box/params"
	issueutils "github.com/hashgard/hashgard/x/issue/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box/types"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/hashgard/hashgard/x/box"
)

func TestDepositBoxList(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})
	cap := 1000
	for i := 0; i < cap; i++ {
		boxInfo := createDepositBox(t, ctx, keeper)
		_, err := keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.Coins{boxInfo.TotalAmount.Token})
		require.Nil(t, err)
	}

	id := ""
	for i := 0; i < 100; i++ {
		//fmt.Println("==================page:" + strconv.Itoa(i))
		boxs := keeper.List(ctx, params.BoxQueryParams{StartId: id, BoxType: types.Deposit, Owner: nil, Limit: 10})
		require.Len(t, boxs, 10)
		for j, box := range boxs {

			if j > 0 {
				require.True(t, boxs[j].CreatedTime <= (boxs[j-1].CreatedTime))
			}
			//fmt.Println(box.Id + "----" + box.CreatedTime.String())
			id = box.Id
		}

	}

}
func TestQueryDepositListFromDepositBox(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	_, err := keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, boxInfo.Owner, boxInfo.Deposit.Interest.Token)
	require.Nil(t, err)

	boxInfo = keeper.GetBox(ctx, boxInfo.Id)
	err = keeper.ProcessDepositBoxByEndBlocker(ctx, boxInfo)
	require.Nil(t, err)

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))
	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.TotalAmount.Token))

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		issueutils.MulDecimals(sdk.NewInt(5000), TestTokenDecimals)), types.Inject)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)), types.Inject)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)), types.Inject)
	require.Nil(t, err)

}
