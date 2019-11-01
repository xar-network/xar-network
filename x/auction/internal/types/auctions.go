package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// MaxAuctionDuration max length of auction, in blocks
	MaxAuctionDuration EndTime = 2 * 24 * 3600 / 5 // roughly 2 days, at 5s block time // 34560
	// BidDuration how long an auction gets extended when someone bids, in blocks
	BidDuration EndTime = 3 * 3600 / 5 // roughly 3 hours, at 5s block time TODO better name // 2160
)

// Auction is an interface to several types of auction.
type Auction interface {
	GetID() ID
	SetID(ID)
	PlaceBid(currentBlockHeight EndTime, bidder sdk.AccAddress, lot sdk.Coin, bid sdk.Coin) ([]BankOutput, []bankInput, sdk.Error)
	GetEndTime() EndTime // auctions close at the end of the block with blockheight EndTime (ie bids placed in that block are valid)
	GetPayout() bankInput
	String() string
}

// BaseAuction type shared by all Auctions
type BaseAuction struct {
	ID         ID
	Initiator  sdk.AccAddress // Person who starts the auction. Giving away Lot (aka seller in a forward auction)
	Lot        sdk.Coin       // Amount of coins up being given by initiator (FA - amount for sale by seller, RA - cost of good by buyer (bid))
	Bidder     sdk.AccAddress // Person who bids in the auction. Receiver of Lot. (aka buyer in forward auction, seller in RA)
	Bid        sdk.Coin       // Amount of coins being given by the bidder (FA - bid, RA - amount being sold)
	EndTime    EndTime        // Block height at which the auction closes. It closes at the end of this block
	MaxEndTime EndTime        // Maximum closing time. Auctions can close before this but never after.
}

// ID type for auction IDs
type ID uint64

// NewIDFromString generate new auction ID from a string
func NewIDFromString(s string) (ID, error) {
	n, err := strconv.ParseUint(s, 10, 64) // copied from how the gov module rest handler's parse proposal IDs
	if err != nil {
		return 0, err
	}
	return ID(n), nil
}

type EndTime int64 // TODO rename to Blockheight or don't define custom type

// Initially the input and output types from the bank module where used here. But they use sdk.Coins instad of sdk.Coin. So it caused a lot of type conversion as auction mainly uses sdk.Coin.
type bankInput struct {
	Address sdk.AccAddress
	Coin    sdk.Coin
}
type BankOutput struct {
	Address sdk.AccAddress
	Coin    sdk.Coin
}

// GetID getter for auction ID
func (a BaseAuction) GetID() ID { return a.ID }

// SetID setter for auction ID
func (a *BaseAuction) SetID(id ID) { a.ID = id }

// GetEndTime getter for auction end time
func (a BaseAuction) GetEndTime() EndTime { return a.EndTime }

// GetPayout implements Auction
func (a BaseAuction) GetPayout() bankInput {
	return bankInput{a.Bidder, a.Lot}
}

func (e EndTime) String() string {
	return string(e)
}

func (a BaseAuction) String() string {
	return fmt.Sprintf(`Auction %d:
  Initiator:              %s
  Lot:               			%s
  Bidder:            		  %s
  Bid:        						%s
  End Time:   						%s
  Max End Time:      			%s`,
		a.GetID(), a.Initiator, a.Lot,
		a.Bidder, a.Bid, a.GetEndTime().String(),
		a.MaxEndTime.String(),
	)
}

// ForwardAuction type for forward auctions
type ForwardAuction struct {
	BaseAuction
}

// NewForwardAuction creates a new forward auction
func NewForwardAuction(seller sdk.AccAddress, lot sdk.Coin, initialBid sdk.Coin, endTime EndTime) (ForwardAuction, BankOutput) {
	auction := ForwardAuction{BaseAuction{
		// no ID
		Initiator:  seller,
		Lot:        lot,
		Bidder:     seller,     // send the proceeds from the first bid back to the seller
		Bid:        initialBid, // set this to zero most of the time
		EndTime:    endTime,
		MaxEndTime: endTime,
	}}
	output := BankOutput{seller, lot}
	return auction, output
}

// PlaceBid implements Auction
func (a *ForwardAuction) PlaceBid(currentBlockHeight EndTime, bidder sdk.AccAddress, lot sdk.Coin, bid sdk.Coin) ([]BankOutput, []bankInput, sdk.Error) {
	// TODO check lot size matches lot?
	// check auction has not closed
	if currentBlockHeight > a.EndTime {
		return []BankOutput{}, []bankInput{}, sdk.ErrInternal("auction has closed")
	}
	// check bid is greater than last bid
	if !a.Bid.IsLT(bid) { // TODO add minimum bid size
		return []BankOutput{}, []bankInput{}, sdk.ErrInternal("bid not greater than last bid")
	}
	// calculate coin movements
	outputs := []BankOutput{{bidder, bid}}                                  // new bidder pays bid now
	inputs := []bankInput{{a.Bidder, a.Bid}, {a.Initiator, bid.Sub(a.Bid)}} // old bidder is paid back, extra goes to seller

	// update auction
	a.Bidder = bidder
	a.Bid = bid
	// increment timeout // TODO into keeper?
	a.EndTime = EndTime(min(int64(currentBlockHeight+BidDuration), int64(a.MaxEndTime))) // TODO is there a better way to structure these types?

	return outputs, inputs, nil
}

// ReverseAuction type for reverse auctions
type ReverseAuction struct {
	BaseAuction
}

// NewReverseAuction creates a new reverse auction
func NewReverseAuction(buyer sdk.AccAddress, bid sdk.Coin, initialLot sdk.Coin, endTime EndTime) (ReverseAuction, BankOutput) {
	auction := ReverseAuction{BaseAuction{
		// no ID
		Initiator:  buyer,
		Lot:        initialLot,
		Bidder:     buyer, // send proceeds from the first bid to the buyer
		Bid:        bid,   // amount that the buyer it buying - doesn't change over course of auction
		EndTime:    endTime,
		MaxEndTime: endTime,
	}}
	output := BankOutput{buyer, initialLot}
	return auction, output
}

// PlaceBid implements Auction
func (a *ReverseAuction) PlaceBid(currentBlockHeight EndTime, bidder sdk.AccAddress, lot sdk.Coin, bid sdk.Coin) ([]BankOutput, []bankInput, sdk.Error) {

	// check bid size matches bid?
	// check auction has not closed
	if currentBlockHeight > a.EndTime {
		return []BankOutput{}, []bankInput{}, sdk.ErrInternal("auction has closed")
	}
	// check bid is less than last bid
	if !lot.IsLT(a.Lot) { // TODO add min bid decrements
		return []BankOutput{}, []bankInput{}, sdk.ErrInternal("lot not smaller than last lot")
	}
	// calculate coin movements
	outputs := []BankOutput{{bidder, a.Bid}}                                // new bidder pays bid now
	inputs := []bankInput{{a.Bidder, a.Bid}, {a.Initiator, a.Lot.Sub(lot)}} // old bidder is paid back, decrease in price for goes to buyer

	// update auction
	a.Bidder = bidder
	a.Lot = lot
	// increment timeout // TODO into keeper?
	a.EndTime = EndTime(min(int64(currentBlockHeight+BidDuration), int64(a.MaxEndTime))) // TODO is there a better way to structure these types?

	return outputs, inputs, nil
}

// ForwardReverseAuction type for forward reverse auction
type ForwardReverseAuction struct {
	BaseAuction
	MaxBid      sdk.Coin
	OtherPerson sdk.AccAddress // TODO rename, this is normally the original CSDT owner
}

func (a ForwardReverseAuction) String() string {
	return fmt.Sprintf(`Auction %d:
  Initiator:              %s
  Lot:               			%s
  Bidder:            		  %s
  Bid:        						%s
  End Time:   						%s
	Max End Time:      			%s
	Max Bid									%s
	Other Person						%s`,
		a.GetID(), a.Initiator, a.Lot,
		a.Bidder, a.Bid, a.GetEndTime().String(),
		a.MaxEndTime.String(), a.MaxBid, a.OtherPerson,
	)
}

// NewForwardReverseAuction creates a new forward reverse auction
func NewForwardReverseAuction(seller sdk.AccAddress, lot sdk.Coin, initialBid sdk.Coin, endTime EndTime, maxBid sdk.Coin, otherPerson sdk.AccAddress) (ForwardReverseAuction, BankOutput) {
	auction := ForwardReverseAuction{
		BaseAuction: BaseAuction{
			// no ID
			Initiator:  seller,
			Lot:        lot,
			Bidder:     seller,     // send the proceeds from the first bid back to the seller
			Bid:        initialBid, // 0 most of the time
			EndTime:    endTime,
			MaxEndTime: endTime},
		MaxBid:      maxBid,
		OtherPerson: otherPerson,
	}
	output := BankOutput{seller, lot}
	return auction, output
}

// PlaceBid implements auction
func (a *ForwardReverseAuction) PlaceBid(currentBlockHeight EndTime, bidder sdk.AccAddress, lot sdk.Coin, bid sdk.Coin) (outputs []BankOutput, inputs []bankInput, err sdk.Error) {
	// check auction has not closed
	if currentBlockHeight > a.EndTime {
		return []BankOutput{}, []bankInput{}, sdk.ErrInternal("auction has closed")
	}

	// determine phase of auction
	switch {
	case a.Bid.IsLT(a.MaxBid) && bid.IsLT(a.MaxBid):
		// Forward auction phase
		if !a.Bid.IsLT(bid) { // TODO add min bid increments
			return []BankOutput{}, []bankInput{}, sdk.ErrInternal("bid not greater than last bid")
		}
		outputs = []BankOutput{{bidder, bid}}                                  // new bidder pays bid now
		inputs = []bankInput{{a.Bidder, a.Bid}, {a.Initiator, bid.Sub(a.Bid)}} // old bidder is paid back, extra goes to seller
	case a.Bid.IsLT(a.MaxBid):
		// Switch over phase
		if !bid.IsEqual(a.MaxBid) { // require bid == a.MaxBid
			return []BankOutput{}, []bankInput{}, sdk.ErrInternal("bid greater than the max bid")
		}
		outputs = []BankOutput{{bidder, bid}} // new bidder pays bid now
		inputs = []bankInput{
			{a.Bidder, a.Bid},               // old bidder is paid back
			{a.Initiator, bid.Sub(a.Bid)},   // extra goes to seller
			{a.OtherPerson, a.Lot.Sub(lot)}, //decrease in price for goes to original CSDT owner
		}

	case a.Bid.IsEqual(a.MaxBid):
		// Reverse auction phase
		if !lot.IsLT(a.Lot) { // TODO add min bid decrements
			return []BankOutput{}, []bankInput{}, sdk.ErrInternal("lot not smaller than last lot")
		}
		outputs = []BankOutput{{bidder, a.Bid}}                                  // new bidder pays bid now
		inputs = []bankInput{{a.Bidder, a.Bid}, {a.OtherPerson, a.Lot.Sub(lot)}} // old bidder is paid back, decrease in price for goes to original CSDT owner
	default:
		panic("should never be reached") // TODO
	}

	// update auction
	a.Bidder = bidder
	a.Lot = lot
	a.Bid = bid
	// increment timeout
	a.EndTime = EndTime(min(int64(currentBlockHeight+BidDuration), int64(a.MaxEndTime))) // TODO is there a better way to structure these types?

	return outputs, inputs, nil
}
