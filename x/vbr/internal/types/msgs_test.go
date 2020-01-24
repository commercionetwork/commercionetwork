package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var TestFunder, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestAmount = sdk.NewCoins(sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(100),
})

var msgIncrementsBRPool = MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: TestAmount,
}

var msgIncrementsBrPoolNoFunds = MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: sdk.NewCoins(sdk.Coin{
		Denom:  "ucommercio",
		Amount: sdk.NewInt(0),
	}),
}

func TestMsgIncrementBlockRewardsPool_Route(t *testing.T) {
	actual := msgIncrementsBRPool.Route()
	expected := ModuleName

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_Type(t *testing.T) {
	actual := msgIncrementsBRPool.Type()
	expected := MsgTypeIncrementBlockRewardsPool

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_valid(t *testing.T) {
	actual := msgIncrementsBRPool.ValidateBasic()

	require.Nil(t, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_noFunds(t *testing.T) {
	actual := msgIncrementsBrPoolNoFunds.ValidateBasic()
	expected := sdk.ErrUnknownRequest("You can't transfer a null or negative amount")

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSignBytes(t *testing.T) {
	actual := msgIncrementsBRPool.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgIncrementsBRPool))

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSigners(t *testing.T) {
	actual := msgIncrementsBRPool.GetSigners()
	expected := []sdk.AccAddress{msgIncrementsBRPool.Funder}

	require.Equal(t, expected, actual)
}
