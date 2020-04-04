package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// -----------
// --- Cdp
// -----------

var testOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testCdp = NewCdp(
	testOwner,
	sdk.NewCoin("ucommercio", sdk.NewInt(100)),
	sdk.NewCoins(sdk.NewCoin("ucc", sdk.NewInt(50))),
	10,
)

func TestCdp_Validate(t *testing.T) {
	testData := []struct {
		name          string
		cdp           Cdp
		shouldBeValid bool
		error         error
	}{
		{
			name:          "Invalid CDP owner",
			cdp:           NewCdp(sdk.AccAddress{}, testCdp.Deposit, testCdp.Credits, testCdp.CreatedAt),
			shouldBeValid: false,
			error:         fmt.Errorf("invalid owner address: %s", sdk.AccAddress{}),
		},
		{
			name:          "Invalid deposited amount",
			cdp:           NewCdp(testCdp.Owner, sdk.Coin{}, testCdp.Credits, testCdp.CreatedAt),
			shouldBeValid: false,
			error:         fmt.Errorf("invalid deposit amount: <nil>"),
		},
		{
			name:          "Invalid liquidity amount",
			cdp:           NewCdp(testCdp.Owner, testCdp.Deposit, sdk.Coins{}, testCdp.CreatedAt),
			shouldBeValid: false,
			error:         fmt.Errorf("invalid liquidity amount: %s", sdk.Coins{}),
		},
		{
			name:          "Invalid timestamp",
			cdp:           NewCdp(testCdp.Owner, testCdp.Deposit, testCdp.Credits, 0),
			shouldBeValid: false,
			error:         fmt.Errorf("invalid timestamp: %d", 0),
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldBeValid {
				require.NoError(t, test.cdp.Validate())
			} else {
				res := test.cdp.Validate()
				require.NotNil(t, res)
				require.Equal(t, test.error, res)
			}
		})
	}
}

func TestCdp_Equals(t *testing.T) {
	testData := []struct {
		name          string
		first         Cdp
		second        Cdp
		shouldBeEqual bool
	}{
		{
			name:          "CDPs are identical",
			first:         testCdp,
			second:        testCdp,
			shouldBeEqual: true,
		},
		{
			name:  "CDPs are different",
			first: testCdp,
			second: Cdp{
				Owner:     testCdp.Owner,
				Deposit:   testCdp.Deposit,
				Credits:   testCdp.Credits,
				CreatedAt: testCdp.CreatedAt + 1,
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
