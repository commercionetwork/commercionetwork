package types

import (
	"fmt"
	"testing"

	//"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var TestFunder, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestAmount = sdk.NewCoins(sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(100),
})

var msgIncrementsBRPool = MsgIncrementBlockRewardsPool{
	Funder: TestFunder.String(),
	Amount: TestAmount,
}

var msgIncrementsBrPoolNoFunds = MsgIncrementBlockRewardsPool{
	Funder: TestFunder.String(),
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
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, "You can't transfer a null or negative amount")

	require.Equal(t, expected.Error(), actual.Error())
}

func TestMsgIncrementBlockRewardsPool_GetSignBytes(t *testing.T) {
	actual := msgIncrementsBRPool.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgIncrementsBRPool))

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSigners(t *testing.T) {
	actual := msgIncrementsBRPool.GetSigners()
	funderAddr, _ := sdk.AccAddressFromBech32(msgIncrementsBRPool.Funder)
	expected := []sdk.AccAddress{funderAddr}

	require.Equal(t, expected, actual)
}
// -------------------------
// --- MsgSetVbrParams
// -------------------------
var TestGov, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestEarnRate = sdk.NewDecWithPrec(5,1)

var msgSetVbrParams = MsgSetVbrParams{
	Government: TestGov.String(),
	DistrEpochIdentifier: EpochDay,
	EarnRate: TestEarnRate,
}

var msgSetVbrParamsNoEpochIdentifier = MsgSetVbrParams{
	Government: TestGov.String(),
	DistrEpochIdentifier: "",
	EarnRate: TestEarnRate,
}

func TestMsgSetVbrParams_Route(t *testing.T) {
	actual := msgSetVbrParams.Route()
	expected := ModuleName

	require.Equal(t, expected, actual)
}

func TestMsgSetVbrParams_Type(t *testing.T) {
	actual := msgSetVbrParams.Type()
	expected := MsgTypeSetVbrParams

	require.Equal(t, expected, actual)
}

func TestMsgSetVbrParams_valid(t *testing.T) {
	actual := msgSetVbrParams.ValidateBasic()

	require.Nil(t, actual)
}

func TestMsgSetVbrParams_ValidateBasic_noEpochIdentifier(t *testing.T) {
	actual := msgSetVbrParamsNoEpochIdentifier.ValidateBasic()
	expected := sdkErr.Wrap(sdkErr.ErrInvalidType, fmt.Sprintf("invalid epoch identifier: %s", msgSetVbrParamsNoEpochIdentifier.DistrEpochIdentifier))

	require.Equal(t, expected.Error(), actual.Error())
}

func TestMsgSetVbrParams_GetSignBytes(t *testing.T) {
	actual := msgSetVbrParams.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgSetVbrParams))

	require.Equal(t, expected, actual)
}

func TestMsgSetVbrParams_GetSigners(t *testing.T) {
	actual := msgSetVbrParams.GetSigners()
	govAddr, _ := sdk.AccAddressFromBech32(msgSetVbrParams.Government)
	expected := []sdk.AccAddress{govAddr}

	require.Equal(t, expected, actual)
}