package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// -----------
// --- Position
// -----------
var testCreatedAt, _ = time.Parse("2006", "2006")
var testOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var testEtp = NewPosition(
	testOwner,
	sdk.NewInt(100),
	sdk.NewCoin("uccc", sdk.NewInt(50)),
	"E95613F1-8407-4B28-9B66-25AB5F4A5FD9",
	testCreatedAt,
	sdk.DecProto{Dec: sdk.NewDec(2)},
)

func TestPosition_Validate(t *testing.T) {
	testData := []struct {
		name          string
		etp           Position
		shouldBeValid bool
	}{
		{
			name:          "Invalid etp owner",
			etp:           NewPosition(sdk.AccAddress{}, sdk.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, testCreatedAt, *testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid collateral amount",
			etp:           NewPosition(sdk.AccAddress(testEtp.Owner), sdk.ZeroInt(), *testEtp.Credits, testEtp.ID, testCreatedAt, *testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid liquidity amount",
			etp:           NewPosition(sdk.AccAddress(testEtp.Owner), sdk.NewInt(testEtp.Collateral), sdk.Coin{}, testEtp.ID, testCreatedAt, *testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid timestamp",
			etp:           NewPosition(sdk.AccAddress(testEtp.Owner), sdk.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, time.Time{}, *testEtp.ExchangeRate),
			shouldBeValid: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldBeValid {
				require.NoError(t, test.etp.Validate())
			} else {
				res := test.etp.Validate()
				require.Error(t, res)
			}
		})
	}
}

func TestPosition_Equals(t *testing.T) {
	secondCreatedAt := testCreatedAt.AddDate(0, 0, 1)

	testData := []struct {
		name          string
		first         Position
		second        Position
		shouldBeEqual bool
	}{
		{
			name:          "etps are identical",
			first:         testEtp,
			second:        testEtp,
			shouldBeEqual: true,
		},
		{
			name:  "etps are different",
			first: testEtp,
			second: Position{
				Owner:      testEtp.Owner,
				Collateral: testEtp.Collateral,
				Credits:    testEtp.Credits,
				CreatedAt:  secondCreatedAt.String(), // TODO check correct test
			},
			shouldBeEqual: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.shouldBeEqual, test.first.Equals(test.second))
		})
	}
}
