package types

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var addr, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestFunder = Funder{Address: addr}
var TestAmount = sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(10),
}

var msgIncrementsBRPool = MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: TestAmount,
}

var msgIncrementsBrPoolNoFunds = MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: types.Coin{Amount: sdk.NewInt(0), Denom: DefaultBondDenom},
}

var msgIncrementsBrPoolWrongDenom = MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: types.Coin{Amount: sdk.NewInt(1), Denom: "dogecoin"},
}

func TestMsgIncrementBlockRewardsPool_Route(t *testing.T) {
	actual := msgIncrementsBRPool.Route()
	expected := ModuleName

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_Type(t *testing.T) {
	actual := msgIncrementsBRPool.Type()
	expected := MsgTypeIncrementBlockRewardsPool

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_valid(t *testing.T) {
	actual := msgIncrementsBRPool.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_noFunds(t *testing.T) {
	actual := msgIncrementsBrPoolNoFunds.ValidateBasic()
	expected := sdk.ErrUnknownRequest("You can't transfer a null amount")

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_WrongCoinDenom(t *testing.T) {
	actual := msgIncrementsBrPoolWrongDenom.ValidateBasic()
	expected := sdk.ErrUnknownRequest(fmt.Sprintf("You can't transfer others than %s", DefaultBondDenom))

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSignBytes(t *testing.T) {
	actual := msgIncrementsBRPool.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgIncrementsBRPool))

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSigners(t *testing.T) {
	actual := msgIncrementsBRPool.GetSigners()
	expected := []sdk.AccAddress{msgIncrementsBRPool.Funder.Address}

	assert.Equal(t, expected, actual)
}
