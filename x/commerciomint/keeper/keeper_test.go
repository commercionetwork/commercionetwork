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

func TestKeeper_SetPosition(t *testing.T) {
	ctx, bk, _, _, k := SetupTestInput()
	// handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName),
		sdk.NewCoins(sdk.NewCoin("ucommercio", testEtp.Collateral)),
	)
	_ = bk.SetCoins(ctx, testEtp.Owner, sdk.NewCoins(testEtp.Credits))
	require.Equal(t, 0, len(k.GetAllPositions(ctx)))
	k.SetPosition(ctx, testEtp)
	require.Equal(t, 1, len(k.GetAllPositions(ctx)))
	position, found := k.GetPosition(ctx, testEtp.Owner, testEtp.ID)
	require.True(t, found)
	require.Equal(t, testEtp.Owner, position.Owner)
	require.True(t, testEtp.CreatedAt.Equal(position.CreatedAt))
}

// --------------
// --- etps
// --------------

func TestKeeper_UpdatePositionBasic(t *testing.T) {
	testData := []struct {
		name            string
		position        types.Position
		insPostion      bool
		shouldBeUpdated bool
	}{
		{
			name:            "Etp doesn't exists",
			position:        fakeEtp,
			insPostion:      false,
			shouldBeUpdated: false,
		},

		{
			name:            "Etp update properly",
			position:        testEtp,
			insPostion:      true,
			shouldBeUpdated: true,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, _, k := SetupTestInput()
			if test.insPostion {
				k.SetPosition(ctx, test.position)
			}
			if test.shouldBeUpdated {
				require.NoError(t, k.UpdatePosition(ctx, test.position))
				return
			}

			if !test.shouldBeUpdated {
				require.Error(t, k.UpdatePosition(ctx, test.position))
				return
			}
		})
	}
}

func TestKeeper_SetPositionBasic(t *testing.T) {
	testData := []struct {
		name             string
		positions        []types.Position
		newPosition      types.Position
		shouldBeInserted bool
	}{
		{
			name:             "New etp is inserted properly",
			positions:        nil,
			newPosition:      testEtp,
			shouldBeInserted: true,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, _, k := SetupTestInput()
			for _, etp := range test.positions {
				k.SetPosition(ctx, etp)
			}

			if test.shouldBeInserted {
				require.NotPanics(t, func() { k.SetPosition(ctx, test.newPosition) })
				require.Len(t, k.GetAllPositions(ctx), len(test.positions)+1)
				return
			}

			if !test.shouldBeInserted {
				require.Panics(t, func() { k.SetPosition(ctx, test.newPosition) })
				require.Len(t, k.GetAllPositions(ctx), len(test.positions))
				return
			}
		})
	}
}

func TestKeeper_NewPosition(t *testing.T) {
	testData := []struct {
		name            string
		owner           sdk.AccAddress
		id              string
		amount          sdk.Int
		userFunds       sdk.Coins
		error           error
		returnedCredits sdk.Coins
	}{
		{
			name:   "invalid deposited amount",
			owner:  testEtp.Owner,
			id:     testEtp.ID,
			amount: sdk.NewInt(0),
			error:  fmt.Errorf("no uccc requested"),
		},
		{
			name:   "Not enough funds inside user wallet",
			amount: testEtp.Collateral,
			owner:  testEtp.Owner,
			id:     testEtp.ID,
			error: fmt.Errorf("insufficient funds: insufficient account funds;  < %s",
				sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 200)),
			),
		},
		{
			name:            "Successful opening",
			amount:          testEtp.Collateral,
			owner:           testEtp.Owner,
			id:              testEtp.ID,
			userFunds:       sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 200)),
			returnedCredits: sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)),
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(10)

			// Setup
			if !test.userFunds.Empty() {
				_ = bk.SetCoins(ctx, test.owner, test.userFunds)
			}

			err := k.NewPosition(ctx, test.owner, sdk.NewCoins(sdk.NewCoin("uccc", test.amount)), test.id)
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
		ctx, _, _, _, k := SetupTestInput()
		require.Empty(t, k.GetAllPositionsOwnedBy(ctx, testEtpOwner))
	})

	t.Run("Existing list is returned properly", func(t *testing.T) {
		ctx, _, _, _, k := SetupTestInput()
		k.SetPosition(ctx, testEtp)
		for _, pos := range k.GetAllPositionsOwnedBy(ctx, testEtpOwner) {
			pos.Equals(testEtp)
		}
	})
}

func TestKeeper_BurnCCC(t *testing.T) {
	t.Run("Non existing ETP returns error", func(t *testing.T) {
		ctx, _, _, _, k := SetupTestInput()

		err := k.BurnCCC(ctx, testEtp.Owner, "notExists", testEtp.Credits)
		errMsg := fmt.Sprintf("position for user with address %s and id %s does not exist", testEtpOwner, "notExists")
		require.Error(t, err)
		require.Equal(t, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg).Error(), err.Error())
	})

	t.Run("Existing ETP is closed properly", func(t *testing.T) {
		ctx, bk, _, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_, _ = bk.AddCoins(ctx, testEtpOwner, sdk.NewCoins(testEtp.Credits))

		require.NoError(t, k.BurnCCC(ctx, testEtpOwner, testEtp.ID, testEtp.Credits))
		require.Equal(t, testEtp.Collateral, bk.GetCoins(ctx, testEtpOwner).AmountOf("ucommercio"))
	})

	t.Run("Existing ETP return correct balance", func(t *testing.T) {
		ctx, bk, _, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		baseUcccAccount := sdk.NewCoin("uccc", sdk.NewInt(50))
		baseUcommercioAccount := sdk.NewCoin("ucommercio", sdk.NewInt(0))
		_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_, _ = bk.AddCoins(ctx, testEtpOwner, sdk.NewCoins(baseUcommercioAccount, baseUcccAccount))
		require.NoError(t, k.BurnCCC(ctx, testEtpOwner, testEtp.ID, halfCoinSub))
		require.Equal(t, baseUcccAccount.Amount.Sub(halfCoinSub.Amount), bk.GetCoins(ctx, testEtpOwner).AmountOf("uccc"))

		burnAmountDec := sdk.NewDecFromInt(halfCoinSub.Amount)
		collateralAmount := burnAmountDec.Mul(testEtp.ExchangeRate).Ceil().TruncateInt()

		require.Equal(t, collateralAmount, bk.GetCoins(ctx, testEtpOwner).AmountOf("ucommercio"))
	})

	t.Run("Existing ETP can't modify before freeze period passes", func(t *testing.T) {
		ctx, bk, _, _, k := SetupTestInput()
		_ = k.SetFreezePeriod(ctx, 3000000000) // 30 seconds
		k.SetPosition(ctx, testEtp)
		_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_, _ = bk.AddCoins(ctx, testEtpOwner, sdk.NewCoins(testEtp.Credits))

		require.Error(t, k.BurnCCC(ctx, testEtpOwner, testEtp.ID, testEtp.Credits))
	})

}

func TestKeeper_deletePosition(t *testing.T) {
	testData := []struct {
		name              string
		existingPositions []types.Position
		deletedPosition   types.Position
		shouldBeDeleted   bool
	}{
		{
			name:              "Existing etp is deleted",
			existingPositions: []types.Position{testEtp},
			deletedPosition:   testEtp,
			shouldBeDeleted:   true,
		},
		{
			name:              "Non existent etp is not deleted",
			existingPositions: []types.Position{testEtp},
			deletedPosition: types.Position{
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
			ctx, _, _, _, k := SetupTestInput()

			for _, etp := range test.existingPositions {
				k.SetPosition(ctx, etp)
			}

			if test.shouldBeDeleted {
				require.NotPanics(t, func() { k.deletePosition(ctx, test.deletedPosition) })
			} else {
				require.Panics(t, func() { k.deletePosition(ctx, test.deletedPosition) })
			}

			result := k.GetAllPositions(ctx)
			if test.shouldBeDeleted {
				require.Len(t, result, len(test.existingPositions)-1)
			} else {
				require.Len(t, result, len(test.existingPositions))
			}
		})
	}
}

// --------------
// --- ConversionRate
// --------------

func TestKeeper_SetConversionRate(t *testing.T) {
	ctx, _, _, _, k := SetupTestInput()
	require.Error(t, k.SetConversionRate(ctx, sdk.NewDec(0)))
	require.Error(t, k.SetConversionRate(ctx, sdk.NewDec(-1)))
	require.NoError(t, k.SetConversionRate(ctx, sdk.NewDec(2)))
	rate := sdk.NewDec(3)
	require.NoError(t, k.SetConversionRate(ctx, rate))

	var got sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get([]byte(types.CollateralRateKey)), &got)
	require.True(t, rate.Equal(got), got.String())
}

func TestKeeper_GetConversionRate(t *testing.T) {
	ctx, _, _, _, k := SetupTestInput()
	rate := sdk.NewDec(3)
	require.NoError(t, k.SetConversionRate(ctx, rate))
	require.Equal(t, rate, k.GetConversionRate(ctx))
}
