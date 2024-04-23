package types

import (
	"testing"
	"time"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------
// --- Position
// -----------
var testCreatedAt, _ = time.Parse("2006", "2006")
var testOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var testEtp = NewPosition(
	testOwner,
	math.NewInt(100),
	sdk.NewCoin(CreditsDenom, math.NewInt(50)),
	"E95613F1-8407-4B28-9B66-25AB5F4A5FD9",
	testCreatedAt,
	math.LegacyNewDec(2),
)
var ownerAddr, _ = sdk.AccAddressFromBech32(testEtp.Owner)

func TestPosition_Validate(t *testing.T) {
	testData := []struct {
		name          string
		etp           Position
		shouldBeValid bool
	}{
		{
			name:          "Invalid etp owner",
			etp:           NewPosition(sdk.AccAddress{}, math.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, testCreatedAt, testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid collateral amount",
			etp:           NewPosition(ownerAddr, math.ZeroInt(), *testEtp.Credits, testEtp.ID, testCreatedAt, testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid liquidity amount",
			etp:           NewPosition(ownerAddr, math.NewInt(testEtp.Collateral), sdk.Coin{}, testEtp.ID, testCreatedAt, testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid timestamp",
			etp:           NewPosition(ownerAddr, math.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, time.Time{}, testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "Invalid exchange rate",
			etp:           NewPosition(ownerAddr, math.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, testCreatedAt, math.LegacyNewDec(-1)),
			shouldBeValid: false,
		},
		{
			name:          "Invalid id",
			etp:           NewPosition(ownerAddr, math.NewInt(testEtp.Collateral), *testEtp.Credits, "abcd", testCreatedAt, testEtp.ExchangeRate),
			shouldBeValid: false,
		},
		{
			name:          "ok",
			etp:           NewPosition(ownerAddr, math.NewInt(testEtp.Collateral), *testEtp.Credits, testEtp.ID, testCreatedAt, testEtp.ExchangeRate),
			shouldBeValid: true,
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

	testData := []struct {
		name          string
		etp           func() Position
		shouldBeEqual bool
	}{
		{
			name:          "identical",
			etp:           func() Position { return testEtp },
			shouldBeEqual: true,
		},
		{
			name: "different Collateral",
			etp: func() Position {
				etp := testEtp
				etp.Collateral = math.NewInt(150).ToLegacyDec().RoundInt64()
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different CreatedAt",
			etp: func() Position {
				etp := testEtp
				t := testEtp.CreatedAt.Add(time.Second)
				etp.CreatedAt = &t
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different Credits",
			etp: func() Position {
				etp := testEtp
				assert.False(t, testEtp.Credits.IsZero())
				credits := testEtp.Credits.Add(*testEtp.Credits)
				etp.Credits = &credits
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different ExchangeRate",
			etp: func() Position {
				etp := testEtp
				assert.False(t, testEtp.ExchangeRate.IsZero())
				etp.ExchangeRate = testEtp.ExchangeRate.Add(testEtp.ExchangeRate)
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different ID",
			etp: func() Position {
				etp := testEtp
				etp.ID = testEtp.ID + "A"
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different Owner",
			etp: func() Position {
				etp := testEtp
				etp.Owner = ""
				return etp
			},
			shouldBeEqual: false,
		},
		{
			name: "different Owner",
			etp: func() Position {
				etp := testEtp
				etp.Owner = ""
				return etp
			},
			shouldBeEqual: false,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.shouldBeEqual, testEtp.Equals(tt.etp()))
		})
	}
}

var validDepositCoin = sdk.NewCoin(CreditsDenom, math.NewInt(50))
var inValidDenomDepositCoin = sdk.NewCoin(BondDenom, math.NewInt(10))

func TestValidateDeposit(t *testing.T) {
	type args struct {
		deposit sdk.Coins
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				deposit: []sdk.Coin{validDepositCoin},
			},
			want: true,
		},
		{
			name: "empty",
			args: args{
				deposit: []sdk.Coin{},
			},
		},
		{
			name: "contains invalid coin",
			args: args{
				deposit: []sdk.Coin{validDepositCoin, inValidDenomDepositCoin},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDeposit(tt.args.deposit); got != tt.want {
				t.Errorf("ValidateDeposit() = %v, want %v", got, tt.want)
			}
		})
	}
}
