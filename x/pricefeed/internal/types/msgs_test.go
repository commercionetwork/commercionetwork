package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Test variables
var testOracle, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testPrice = Price{AssetName: "uatom", Value: sdk.NewDecWithPrec(15423, 2), Expiry: sdk.NewInt(1100)}

// -------------------
// --- MsgSetPrice
// -------------------

var msgSetPrice = NewMsgSetPrice(testPrice, testOracle)

func TestMsgSetPrice_Route(t *testing.T) {
	require.Equal(t, RouterKey, msgSetPrice.Route())
}

func TestMsgSetPrice_Type(t *testing.T) {
	require.Equal(t, MsgTypeSetPrice, msgSetPrice.Type())
}

func TestMsgSetPrice_ValidateBasic_ValidMessage(t *testing.T) {
	require.Nil(t, msgSetPrice.ValidateBasic())
}

func TestMsgSetPrice_ValidateBasic_InvalidMessage(t *testing.T) {
	msgInvalid := MsgSetPrice{
		Oracle: testOracle,
		Price: Price{
			AssetName: "    ",
			Value:     sdk.NewDec(10),
			Expiry:    sdk.Int{},
		},
	}
	require.Error(t, msgInvalid.ValidateBasic())
}

func TestMsgSetPrice_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetPrice","value":{"oracle":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","price":{"asset_name":"uatom","expiry":"1100","rate":"0","value":"154.230000000000000000"}}}`
	require.Equal(t, expected, string(msgSetPrice.GetSignBytes()))
}

func TestMsgSetPrice_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetPrice.Oracle}
	require.Equal(t, expected, msgSetPrice.GetSigners())
}

func TestMsgSetPrice_ParseJson(t *testing.T) {
	json := `{"type":"commercio/MsgSetPrice","value":{"oracle":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","price":{"asset_name":"uatom","value":"154.23","expiry":"1100"}}}`

	var msg MsgSetPrice
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)
	require.Equal(t, msgSetPrice, msg)
}

// -------------------
// --- MsgAddOracle
// -------------------

var TestGovernment, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var msgAddOracle = NewMsgAddOracle(testOracle, TestGovernment)

func TestMsgAddOracle_Route(t *testing.T) {
	require.Equal(t, RouterKey, msgAddOracle.Route())
}

func TestMsgAddOracle_Type(t *testing.T) {
	require.Equal(t, MsgTypeAddOracle, msgAddOracle.Type())
}

func TestMsgAddOracle_ValidateBasic_ValidMessage(t *testing.T) {
	require.Nil(t, msgAddOracle.ValidateBasic())
}

func TestMsgAddOracle_ValidateBasic_InvalidMessage(t *testing.T) {
	msgInvalid := MsgAddOracle{Signer: nil, Oracle: nil}
	require.Error(t, msgInvalid.ValidateBasic())
}

func TestMsgAddOracle_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddOracle))
	require.Equal(t, expected, msgAddOracle.GetSignBytes())
}

func TestMsgAddOracle_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgAddOracle.Signer}
	require.Equal(t, expected, msgAddOracle.GetSigners())
}
