package compound

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "moneymarket" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateMarket:
			return handleMsgCreateMarket(ctx, keeper, msg)
		case MsgSupplyMarket:
			return handleMsgSupplyMarket(ctx, keeper, msg)
		case MsgBorrowFromMarket:
			return handleMsgBorrowFromMarket(ctx, keeper, msg)
		case MsgRedeemFromMarket:
			return handleMsgRedeemFromMarket(ctx, keeper, msg)
		case MsgRepayToMarket:
			return handleMsgRepayToMarket(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized moneymarket Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create market
func handleMsgCreateMarket(ctx sdk.Context, keeper Keeper, msg MsgCreateMarket) sdk.Result {
	fmt.Println("hello")
	fmt.Println(msg.Name)
	// if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
	// 	return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	// }
	// if keeper.HasMarketOwner(ctx, msg.Name) {
	// 	err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetMarketOwner(ctx, msg.Name), msg.Bid)
	// 	if err != nil {
	// 		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	// 	}
	// } else {
	_, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.InterestRate) // If so, deduct the Bid amount from the sender
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}
	//}

	keeper.SetMarketOwner(ctx, msg.Name, msg.Symbol, msg.Buyer, msg.TokenName, msg.CollateralToken)
	//keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

func handleMsgSupplyMarket(ctx sdk.Context, keeper Keeper, msg MsgSupplyMarket) sdk.Result {
	// if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
	// 	return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	// }
	// if keeper.HasMarketOwner(ctx, msg.Name) {
	// 	err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetMarketOwner(ctx, msg.Name), msg.Bid)
	// 	if err != nil {
	// 		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	// 	}
	// } else {
	_, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Supplier, msg.LendTokens) // If so, deduct the Bid amount from the sender
	if err != nil {
		return sdk.ErrInsufficientCoins("Supplier does not have enough coins").Result()
	}
	//}

	keeper.SupplyMarketPosition(ctx, msg.Supplier, msg.Market, msg.LendTokens)
	//keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

func handleMsgBorrowFromMarket(ctx sdk.Context, keeper Keeper, msg MsgBorrowFromMarket) sdk.Result {

	_, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Supplier, msg.CollateralTokens) // If so, deduct the Bid amount from the sender
	if err != nil {
		return sdk.ErrInsufficientCoins("Market does not have enough coins").Result()
	}

	_, err1 := keeper.coinKeeper.AddCoins(ctx, msg.Supplier, msg.BorrowTokens) // If so, deduct the Bid amount from the sender
	if err1 != nil {
		return sdk.ErrInsufficientCoins("Market does not have enough coins").Result()
	}

	keeper.BorrowFromMarketPosition(ctx, msg.Supplier, msg.Market, msg.BorrowTokens, msg.CollateralTokens)
	//keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

func handleMsgRedeemFromMarket(ctx sdk.Context, keeper Keeper, msg MsgRedeemFromMarket) sdk.Result {
	_, err := keeper.coinKeeper.AddCoins(ctx, msg.Supplier, msg.RedeemTokens) // If so, deduct the Bid amount from the sender
	if err != nil {
		return sdk.ErrInsufficientCoins("Market does not have enough coins").Result()
	}

	keeper.RedeemFromMarketPosition(ctx, msg.Supplier, msg.Market, msg.RedeemTokens)
	//keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

func handleMsgRepayToMarket(ctx sdk.Context, keeper Keeper, msg MsgRepayToMarket) sdk.Result {

	//moneymarket := keeper.GetMarketInfo(ctx, msg.Market)
	marketposition := keeper.GetMarketPosition(ctx, msg.Borrower)
	fmt.Println("market position")
	fmt.Println(marketposition)
	//collateralToReturn := sdk.int64(marketposition.BorrowCollateral[0].Amount)
	//collateralToRemint := sdk.Coins{sdk.NewInt64Coin(moneymarket.CollateralToken, collateralToReturn)}
	_, err := keeper.coinKeeper.AddCoins(ctx, msg.Borrower, marketposition.BorrowCollateral) // If so, deduct the Bid amount from the sender
	if err != nil {
		return sdk.ErrInsufficientCoins("You do not have enough coins").Result()
	}

	_, err1 := keeper.coinKeeper.SubtractCoins(ctx, msg.Borrower, msg.RepayTokens) // If so, deduct the Bid amount from the sender
	if err1 != nil {
		return sdk.ErrInsufficientCoins("You do not have enough coins").Result()
	}

	keeper.RepayToMarketPosition(ctx, msg.Borrower, msg.Market, msg.RepayTokens, marketposition.BorrowCollateral)
	//keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}
