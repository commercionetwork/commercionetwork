package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestDepositedAmount = sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100)))
var TestLiquidityAmount = sdk.NewCoins(sdk.NewCoin("ucc", sdk.NewInt(50)))
var TestOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var timezone, _ = time.LoadLocation("UTC")
var TestTimestamp = time.Date(1990, 01, 01, 20, 20, 00, 0, timezone)

var TestCdp = Cdp{
	Owner:           TestOwner,
	DepositedAmount: TestDepositedAmount,
	CreditsAmount:   TestLiquidityAmount,
	Timestamp:       TestTimestamp,
}

var TestCdpS = Cdps{TestCdp}

func TestCdp_Validate_ValidCdp(t *testing.T) {
	err := TestCdp.Validate()
	assert.Nil(t, err)
}

func TestCdp_Validate_InvalidCdpOwner(t *testing.T) {
	var TestCdp = Cdp{
		Owner:           sdk.AccAddress{},
		DepositedAmount: TestDepositedAmount,
		CreditsAmount:   TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidAddress(TestCdp.Owner.String()), err)
	assert.Error(t, err)
}

func TestCdp_Validate_InvalidCdpDepositedAmount(t *testing.T) {
	var TestCdp = Cdp{
		Owner:           TestOwner,
		DepositedAmount: sdk.Coins{},
		CreditsAmount:   TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidCoins(TestCdp.DepositedAmount.String()), err)
	assert.Error(t, err)
}

func TestCdp_Validate_InvalidCdpLiquidityAmount(t *testing.T) {
	var TestCdp = Cdp{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		CreditsAmount:   sdk.Coins{},
		Timestamp:       TestTimestamp,
	}

	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrInvalidCoins(TestCdp.CreditsAmount.String()), err)
	assert.Error(t, err)
}

func TestCdp_Validate_InvalidCdpTimestamp(t *testing.T) {
	var TestCdp = Cdp{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		CreditsAmount:   TestLiquidityAmount,
		Timestamp:       time.Time{},
	}

	err := TestCdp.Validate()
	assert.Equal(t, sdk.ErrUnknownRequest("timestamp not valid"), err)
	assert.Error(t, err)
}

func TestCdp_Equals_True(t *testing.T) {
	var TestCdp2 = Cdp{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		CreditsAmount:   TestLiquidityAmount,
		Timestamp:       TestTimestamp,
	}
	actual := TestCdp.Equals(TestCdp2)
	assert.True(t, actual)
}

func TestCdp_Equals_false(t *testing.T) {
	var TestCdp2 = Cdp{
		Owner:           TestOwner,
		DepositedAmount: TestDepositedAmount,
		CreditsAmount:   TestLiquidityAmount,
		Timestamp:       time.Time{},
	}
	actual := TestCdp.Equals(TestCdp2)
	assert.False(t, actual)
}

func TestCdps_AppendIfMissing_notMissing(t *testing.T) {
	cdps, found := TestCdpS.AppendIfMissing(TestCdp)
	assert.Nil(t, cdps)
	assert.True(t, found)
}

func TestCdps_AppendIfMissing_Missing(t *testing.T) {
	cdps := Cdps{}
	cdps, found := cdps.AppendIfMissing(TestCdp)
	assert.False(t, found)
	assert.NotNil(t, cdps)
}

func TestCdps_RemoveWhenFound_removed(t *testing.T) {
	cdps, removed := TestCdpS.RemoveWhenFound(TestTimestamp)
	assert.True(t, removed)
	assert.Len(t, cdps, 0)
}

func TestCdps_RemoveWhenFound_notRemoved(t *testing.T) {
	cdps, removed := TestCdpS.RemoveWhenFound(time.Time{})
	assert.False(t, removed)
	assert.Len(t, cdps, 1)
}
