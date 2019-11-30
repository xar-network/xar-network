package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = MsgCreateMarket{}
)

type MsgCreateMarket struct {
	Nominee    sdk.AccAddress `json:"nominee" yaml:"nominee"`
	BaseAsset  string         `json:"base_asset" yaml:"base_asset"`
	QuoteAsset string         `json:"quote_asset" yaml:"quote_asset"`
}

func NewMsgCreateMarket(
	nominee sdk.AccAddress,
	baseAsset string,
	quoteAsset string,
) MsgCreateMarket {
	return MsgCreateMarket{
		Nominee:    nominee,
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

	if msg.Nominee.Empty() {
		return sdk.ErrInvalidAddress("missing nominee address")
	}

	return nil
}

func (msg MsgCreateMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Nominee}
}

func (msg MsgCreateMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
