package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = MsgCreateIssuer{}
	_ sdk.Msg = MsgDestroyIssuer{}
	_ sdk.Msg = MsgCreateOracle{}
)

type (
	MsgCreateIssuer struct {
		Issuer        sdk.AccAddress
		Denominations []string
		Authority     sdk.AccAddress
	}
	MsgDestroyIssuer struct {
		Issuer    sdk.AccAddress
		Authority sdk.AccAddress
	}
	MsgCreateOracle struct {
		Oracle    sdk.AccAddress
		Authority sdk.AccAddress
	}
	MsgCreateMarket struct {
		Authority  sdk.AccAddress
		BaseAsset  string
		QuoteAsset string
	}
)

func (msg MsgDestroyIssuer) Type() string { return "destroyIssuer" }

func (msg MsgCreateIssuer) Type() string { return "createIssuer" }

func (msg MsgCreateOracle) Type() string { return "createOracle" }

func (msg MsgCreateMarket) Type() string { return "createMarket" }

func (msg MsgDestroyIssuer) ValidateBasic() sdk.Error {
	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	return nil
}

func (msg MsgCreateOracle) ValidateBasic() sdk.Error {
	if msg.Oracle.Empty() {
		return sdk.ErrInvalidAddress("missing oracle address")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	return nil
}

func (msg MsgCreateMarket) ValidateBasic() sdk.Error {
	//TODO check if asset exists in supply
	if msg.BaseAsset == "" {
		return sdk.ErrInvalidAddress("missing base asset")
	}

	//TODO check if asset exists in supply
	if msg.QuoteAsset == "" {
		return sdk.ErrInvalidAddress("missing base asset")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	return nil
}

func (msg MsgCreateIssuer) ValidateBasic() sdk.Error {
	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	if len(msg.Denominations) == 0 {
		return ErrNoDenomsSpecified()
	}

	return nil
}

func (msg MsgDestroyIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgCreateIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgCreateOracle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgCreateMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgDestroyIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDestroyIssuer) Route() string { return ModuleName }

func (msg MsgCreateIssuer) Route() string { return ModuleName }

func (msg MsgCreateOracle) Route() string { return ModuleName }

func (msg MsgCreateMarket) Route() string { return ModuleName }
