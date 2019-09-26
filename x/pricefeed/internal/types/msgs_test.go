package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test variables
var testOracle1, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var TestPriceInfo = CurrentPrice{
	AssetName: "test",
	AssetCode: "0000",
	Price:     sdk.NewInt(10),
	Expiry:    sdk.NewInt(5000),
}

var TestRawPrice = RawPrice{
	PriceInfo: TestPriceInfo,
	Oracle:    testOracle1,
}

var msgSetPrice = NewMsgSetPrice(TestRawPrice)

func TestMsgSetPrice_Route(t *testing.T) {
	expected := RouterKey
	actual := msgSetPrice.Route()
	assert.Equal(t, expected, actual)
}

func TestMsgSetPrice_Type(t *testing.T) {
	expected := MsgTypeSetPrice
	actual := msgSetPrice.Type()
	assert.Equal(t, expected, actual)
}

func TestMsgSetPrice_ValidateBasic_ValidMessage(t *testing.T) {
	actual := msgSetPrice.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgSetPrice_ValidateBasic_InvalidMessage(t *testing.T) {
	priceInfo := CurrentPrice{
		AssetName: "    ",
		AssetCode: "1",
		Price:     sdk.NewInt(10),
		Expiry:    sdk.Int{},
	}

	msgInvalid := MsgSetPrice{
		Price: RawPrice{
			Oracle:    testOracle1,
			PriceInfo: priceInfo},
	}

	actual := msgInvalid.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetPrice_GetSignBytes(t *testing.T) {
	actual := msgSetPrice.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSetPrice))
	assert.Equal(t, expected, actual)
}

func TestMsgSetPrice_GetSigners(t *testing.T) {
	actual := msgSetPrice.GetSigners()
	expected := []sdk.AccAddress{msgSetPrice.Price.Oracle}
	assert.Equal(t, expected, actual)
}

///////////////////////////
//////MsgAddOracle////////
/////////////////////////
var TestGovernment, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var msgAddOracle = MsgAddOracle{
	Signer: TestGovernment,
	Oracle: testOracle1,
}

func TestMsgAddOracle_Route(t *testing.T) {
	actual := msgAddOracle.Route()
	expected := RouterKey
	assert.Equal(t, expected, actual)
}

func TestMsgAddOracle_Type(t *testing.T) {
	actual := msgAddOracle.Type()
	expected := MsgTypeAddOracle
	assert.Equal(t, expected, actual)
}

func TestMsgAddOracle_ValidateBasic_ValidMessage(t *testing.T) {
	actual := msgAddOracle.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgAddOracle_ValidateBasic_InvalidMessage(t *testing.T) {
	msgInvalid := MsgAddOracle{
		Signer: nil,
		Oracle: nil,
	}
	actual := msgInvalid.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgAddOracle_GetSignBytes(t *testing.T) {
	actual := msgAddOracle.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddOracle))

	assert.Equal(t, expected, actual)
}

func TestMsgAddOracle_GetSigners(t *testing.T) {
	actual := msgAddOracle.GetSigners()
	expected := []sdk.AccAddress{msgAddOracle.Signer}

	assert.Equal(t, expected, actual)
}
