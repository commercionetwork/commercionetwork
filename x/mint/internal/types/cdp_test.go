package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestDepositedAmount = sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100)))
var TestLiquidityAmount = sdk.NewCoins(sdk.NewCoin("ucc", sdk.NewInt(50)))
var TestOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestTimestamp = "timestamp-test"

var TestCdp = CDP{
	Owner:           TestOwner,
	DepositedAmount: TestDepositedAmount,
	LiquidityAmount: TestLiquidityAmount,
	Timestamp:       TestTimestamp,
}

var TestCdpS = CDPs{TestCdp}

func TestCDP_Validate_ValidCdp(t *testing.T) {
	err := TestCdp.Validate()
	assert.Nil(t, err)
}

func TestCDP_Validate_InvalidCdpOwner(t *testing.T) {
	var TestCdp = CDP{
		Owner:           sdk.AccAddress{},
		DepositedAmount: TestDepositedAmount,
		LiquidityAmount: TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidAddress(TestCdp.Owner.String()), err)
	assert.Error(t, err)
}

func TestCDP_Validate_InvalidCdpDepositedAmount(t *testing.T) {
	var TestCdp = CDP{
		Owner:           TestOwner,
		DepositedAmount: sdk.Coins{},
		LiquidityAmount: TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidCoins(TestCdp.DepositedAmount.String()), err)
	assert.Error(t, err)
}

func TestCDP_Validate_InvalidCdpLiquidityAmount(t *testing.T) {
	var TestCdp = CDP{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		LiquidityAmount: sdk.Coins{},
		Timestamp:       TestTimestamp,
	}

	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidCoins(TestCdp.LiquidityAmount.String()), err)
	assert.Error(t, err)
}

func TestCDP_Validate_InvalidCdpTimestamp(t *testing.T) {
	var TestCdp = CDP{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		LiquidityAmount: TestLiquidityAmount,
		Timestamp:       "    ",
	}

	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrUnknownRequest("timestamp cant be empty"), err)
	assert.Error(t, err)
}

func TestCDP_Equals_True(t *testing.T) {
	var TestCdp2 = CDP{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		LiquidityAmount: TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	actual := TestCdp.Equals(TestCdp2)
	assert.True(t, actual)
}

func TestCDP_Equals_false(t *testing.T) {
	var TestCdp2 = CDP{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		LiquidityAmount: TestLiquidityAmount,
		Timestamp:       "    ",
	}
	actual := TestCdp.Equals(TestCdp2)
	assert.False(t, actual)
}

func TestCDPs_AppendIfMissing_notMissing(t *testing.T) {
	cdps, found := TestCdpS.AppendIfMissing(TestCdp)
	assert.Nil(t, cdps)
	assert.True(t, found)
}

func TestCDPs_AppendIfMissing_Missing(t *testing.T) {
	cdps := CDPs{}
	cdps, found := cdps.AppendIfMissing(TestCdp)
	assert.False(t, found)
	assert.NotNil(t, cdps)
}

func TestCDPs_RemoveWhenFound_removed(t *testing.T) {
	cdps, removed := TestCdpS.RemoveWhenFound(TestTimestamp)
	assert.True(t, removed)
	assert.Len(t, cdps, 0)
}

func TestCDPs_RemoveWhenFound_notRemoved(t *testing.T) {
	cdps, removed := TestCdpS.RemoveWhenFound("tt")
	assert.False(t, removed)
	assert.Len(t, cdps, 1)
}

func TestCDPs_GetCdpFromTimestamp_found(t *testing.T) {
	cdp, found := TestCdpS.GetCdpFromTimestamp(TestTimestamp)
	assert.True(t, found)
	assert.Equal(t, &TestCdp, cdp)
}

func TestCDPs_GetCdpFromTimestamp_notFound(t *testing.T) {
	cdp, found := TestCdpS.GetCdpFromTimestamp("")
	assert.False(t, found)
	assert.Nil(t, cdp)
}
