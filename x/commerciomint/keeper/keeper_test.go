package keeper

import (
	"errors"
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_StoreCdp(t *testing.T) {
	ctx, bk, _, _, _, k := SetupTestInput()
	// handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName),
		sdk.NewCoins(sdk.NewCoin("ucommercio", testEtp.Collateral)),
	)
	_ = bk.SetCoins(ctx, testEtp.Owner, sdk.NewCoins(testEtp.Credits))
	require.Equal(t, 0, len(k.GetAllPositions(ctx)))
	k.SetPosition(ctx, testEtp)
	require.Equal(t, 1, len(k.GetAllPositions(ctx)))
	cdp, found := k.GetPosition(ctx, testEtp.Owner, testEtp.ID)
	require.True(t, found)
	require.Equal(t, testEtp.Owner, cdp.Owner)
	require.True(t, testEtp.CreatedAt.Equal(cdp.CreatedAt))
}

// --------------
// --- CDPs
// --------------

func TestKeeper_StoreCdpBasic(t *testing.T) {
	testData := []struct {
		name             string
		cdps             []types.Position
		newCdp           types.Position
		shouldBeInserted bool
	}{
		{
			name:             "New CDP is inserted properly",
			cdps:             nil,
			newCdp:           testEtp,
			shouldBeInserted: true,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, _, _, k := SetupTestInput()
			for _, cdp := range test.cdps {
				k.SetPosition(ctx, cdp)
			}

			if test.shouldBeInserted {
				require.NotPanics(t, func() { k.SetPosition(ctx, test.newCdp) })
				require.Len(t, k.GetAllPositions(ctx), len(test.cdps)+1)
				return
			}

			if !test.shouldBeInserted {
				require.Panics(t, func() { k.SetPosition(ctx, test.newCdp) })
				require.Len(t, k.GetAllPositions(ctx), len(test.cdps))
				return
			}
		})
	}
}

func TestKeeper_OpenCdp(t *testing.T) {
	testData := []struct {
		name            string
		owner           sdk.AccAddress
		amount          sdk.Int
		userFunds       sdk.Coins
		error           error
		returnedCredits sdk.Coins
	}{
		{
			name:   "invalid deposited amount",
			owner:  testEtp.Owner,
			amount: sdk.NewInt(0),
			error:  fmt.Errorf("no uccc requested"),
		},
		{
			name:   "Not enough funds inside user wallet",
			amount: testEtp.Collateral,
			owner:  testEtp.Owner,
			error: fmt.Errorf("insufficient funds: insufficient account funds;  < %s",
				sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 200)),
			),
		},
		{
			name:            "Successful opening",
			amount:          testEtp.Collateral,
			owner:           testEtp.Owner,
			userFunds:       sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(200))),
			returnedCredits: sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)),
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, _, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(10)

			// Setup
			if !test.userFunds.Empty() {
				_ = bk.SetCoins(ctx, test.owner, test.userFunds)
			}

			err := k.NewPosition(ctx, test.owner, sdk.NewCoins(sdk.NewCoin("uccc", test.amount)))
			if test.error != nil {
				require.Error(t, err)
				e := errors.Unwrap(err)
				if e != nil {
					// error should be unwrapped to match against the exact error string
					err = e
				}
				require.Equal(t, test.error.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if !test.returnedCredits.IsEqual(sdk.Coins{}) {
				actual := bk.GetCoins(ctx, test.owner)
				require.Equal(t, test.returnedCredits, actual)
			}
		})
	}

}

func TestKeeper_GetAllPositionsOwnedBy(t *testing.T) {
	t.Run("Empty list is returned properly", func(t *testing.T) {
		ctx, _, _, _, _, k := SetupTestInput()
		require.Empty(t, k.GetAllPositionsOwnedBy(ctx, testCdpOwner))
	})

	t.Run("Existing list is returned properly", func(t *testing.T) {
		ctx, _, _, _, _, k := SetupTestInput()
		k.SetPosition(ctx, testEtp)
		for _, pos := range k.GetAllPositionsOwnedBy(ctx, testCdpOwner) {
			pos.Equals(testEtp)
		}
	})
}

func TestKeeper_CloseCdp(t *testing.T) {
	t.Run("Non existing CDP returns error", func(t *testing.T) {
		ctx, _, _, _, _, k := SetupTestInput()

		err := k.BurnCCC(ctx, testEtp.Owner, "notExists", testEtp.Credits)
		errMsg := fmt.Sprintf("position for user with address %s and id %s does not exist", testCdpOwner, "notExists")
		require.Equal(t, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg).Error(), err.Error())
	})

	t.Run("Existing CDP is closed properly", func(t *testing.T) {
		ctx, bk, _, _, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_, _ = bk.AddCoins(ctx, testCdpOwner, sdk.NewCoins(testEtp.Credits))

		require.NoError(t, k.BurnCCC(ctx, testCdpOwner, testEtp.ID, testEtp.Credits))
		require.Equal(t, testEtp.Collateral, bk.GetCoins(ctx, testCdpOwner).AmountOf("ucommercio"))
	})

}

func TestKeeper_DeleteCdp(t *testing.T) {
	testData := []struct {
		name            string
		existingCdps    []types.Position
		deletedCdp      types.Position
		shouldBeDeleted bool
	}{
		{
			name:            "Existing CDP is deleted",
			existingCdps:    []types.Position{testEtp},
			deletedCdp:      testEtp,
			shouldBeDeleted: true,
		},
		{
			name:         "Non existent CDP is not deleted",
			existingCdps: []types.Position{testEtp},
			deletedCdp: types.Position{
				Owner:      testEtp.Owner,
				Collateral: testEtp.Collateral,
				Credits:    testEtp.Credits,
				CreatedAt:  testEtp.CreatedAt.AddDate(0, 0, 1),
			},
			shouldBeDeleted: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, _, _, k := SetupTestInput()

			for _, cdp := range test.existingCdps {
				k.SetPosition(ctx, cdp)
			}

			if test.shouldBeDeleted {
				require.NotPanics(t, func() { k.deletePosition(ctx, test.deletedCdp) })
			} else {
				require.Panics(t, func() { k.deletePosition(ctx, test.deletedCdp) })
			}

			result := k.GetAllPositions(ctx)
			if test.shouldBeDeleted {
				require.Len(t, result, len(test.existingCdps)-1)
			} else {
				require.Len(t, result, len(test.existingCdps))
			}
		})
	}
}

// --------------
// --- CdpCollateralRate
// --------------

func TestKeeper_SetCdpCollateralRate(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	require.Error(t, k.SetConversionRate(ctx, sdk.NewInt(0)))
	require.Error(t, k.SetConversionRate(ctx, sdk.NewInt(-1)))
	require.NoError(t, k.SetConversionRate(ctx, sdk.NewInt(2)))
	rate := sdk.NewInt(3)
	require.NoError(t, k.SetConversionRate(ctx, rate))

	var got sdk.Int
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get([]byte(types.CollateralRateKey)), &got)
	require.True(t, rate.Equal(got), got.String())
}

func TestKeeper_GetCdpCollateralRate(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	rate := sdk.NewInt(3)
	require.NoError(t, k.SetConversionRate(ctx, rate))
	require.Equal(t, rate, k.GetConversionRate(ctx))
}
