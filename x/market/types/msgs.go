package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = MsgCreateMarket{}
)

type MsgCreateMarket struct {
	POA        sdk.AccAddress
	BaseAsset  string
	QuoteAsset string
}

func NewMsgCreateMarket(
	poa sdk.AccAddress,
	baseAsset string,
	quoteAsset string,
) MsgCreateMarket {
	return MsgCreateMarket{
		POA:        poa,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}
}
func (msg MsgCreateMarket) Route() string { return ModuleName }

func (msg MsgCreateMarket) Type() string { return "createMarket" }

func (msg MsgCreateMarket) ValidateBasic() sdk.Error {
	if msg.BaseAsset == "" {
		return sdk.ErrInvalidAddress("missing base asset")
	}

	//TODO check if asset exists in supply
	if msg.QuoteAsset == "" {
		return sdk.ErrInvalidAddress("missing base asset")
	}

	if msg.POA.Empty() {
		return sdk.ErrInvalidAddress("missing POA address")
	}

	return nil
}

func (msg MsgCreateMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.POA}
}

func (msg MsgCreateMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
