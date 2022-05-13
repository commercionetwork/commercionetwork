package keeper

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SetPosition(t *testing.T) {
	ctx, bk, _, k := SetupTestInput()
	coins := sdk.NewCoins(*testEtp.Credits)
	err := bk.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		require.NoError(t, err)
	}
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, testEtpOwner, coins)
	if err != nil {
		require.NoError(t, err)
	}
	/*err := bk.SetBalance(ctx, testEtpOwner, *testEtp.Credits)
	if err != nil {
		require.NoError(t, err)
	}*/
	require.Equal(t, 0, len(k.GetAllPositions(ctx)))
	k.SetPosition(ctx, testEtp)
	require.Equal(t, 1, len(k.GetAllPositions(ctx)))
	position, found := k.GetPosition(ctx, testEtpOwner, testEtp.ID)
	require.True(t, found)
	require.Equal(t, testEtp.Owner, position.Owner)
	require.True(t, testEtp.CreatedAt.Equal(*position.CreatedAt))

	// a position with id already exists
	err = k.SetPosition(ctx, testEtp)
	require.Error(t, err)

	invalidTestEtp := testEtp
	invalidTestEtp.Owner = ""
	err = k.SetPosition(ctx, invalidTestEtp)
	require.Error(t, err)
}

// --------------
// --- etps
// --------------

func TestKeeper_UpdatePosition(t *testing.T) {
	testData := []struct {
		name            string
		position        func() types.Position
		insPosition     bool
		shouldBeUpdated bool
	}{
		{
			name: "invalid owner",
			position: func() types.Position {
				pos := testEtp
				pos.Owner = ""
				return pos
			},
			insPosition:     false,
			shouldBeUpdated: false,
		},
		{
			name:            "Etp doesn't exists",
			position:        func() types.Position { return testEtp },
			insPosition:     false,
			shouldBeUpdated: false,
		},

		{
			name:            "Etp updated properly",
			position:        func() types.Position { return testEtp },
			insPosition:     true,
			shouldBeUpdated: true,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			if test.insPosition {
				require.NoError(t, k.SetPosition(ctx, test.position()))
			}
			if test.shouldBeUpdated {
				require.NoError(t, k.UpdatePosition(ctx, test.position()))
				return
			} else {
				require.Error(t, k.UpdatePosition(ctx, test.position()))
				return
			}
		})
	}
}

func TestKeeper_NewPosition(t *testing.T) {
	testData := []struct {
		name            string
		owner           string
		id              string
		amount          sdk.Int
		userFunds       sdk.Coins
		error           error
		returnedCredits sdk.Coins
	}{
		{
			name:   "invalid owner",
			owner:  "",
			id:     testEtp.ID,
			amount: sdk.NewInt(0),
			error:  fmt.Errorf("empty address string is not allowed"),
		},
		{
			name:   "no uccc requested",
			owner:  testEtpOwner.String(),
			id:     testEtp.ID,
			amount: sdk.NewInt(0),
			error:  fmt.Errorf("no uccc requested"),
		},
		{
			name:   "not enough funds inside user wallet",
			amount: sdk.NewInt(testEtp.Collateral),
			owner:  testEtpOwner.String(),
			id:     testEtp.ID,
			error: fmt.Errorf("0"+types.BondDenom+" is smaller than %s: insufficient funds",
				sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, 200)),
			),
		},
		{
			name:            "ok",
			amount:          sdk.NewInt(testEtp.Collateral),
			owner:           testEtpOwner.String(),
			id:              testEtp.ID,
			userFunds:       sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, 200)),
			returnedCredits: sdk.NewCoins(sdk.NewInt64Coin(types.CreditsDenom, 100)),
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(10)

			// Setup
			if !test.userFunds.Empty() {
				ownerAddr, err := sdk.AccAddressFromBech32(test.owner)
				require.NoError(t, err)
				err = bk.MintCoins(ctx, types.ModuleName, test.userFunds)
				require.NoError(t, err)
				err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, test.userFunds)
				require.NoError(t, err)

				/*err = bk.AddCoins(ctx, ownerAddr, test.userFunds)
				require.NoError(t, err)*/
			}

			err := k.NewPosition(ctx, test.owner, sdk.Coins{sdk.Coin{
				Denom:  types.CreditsDenom,
				Amount: test.amount,
			}}, test.id)
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
				ownerAddr, err := sdk.AccAddressFromBech32(test.owner)
				require.NoError(t, err)
				actual := bk.GetAllBalances(ctx, ownerAddr)
				require.Equal(t, test.returnedCredits, actual)
			}
		})
	}

}

func TestKeeper_GetAllPositionsOwnedBy(t *testing.T) {
	t.Run("Empty list is returned properly", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()
		require.Empty(t, k.GetAllPositionsOwnedBy(ctx, testEtpOwner))
	})

	t.Run("Existing list is returned properly", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()
		k.SetPosition(ctx, testEtp)
		for _, pos := range k.GetAllPositionsOwnedBy(ctx, testEtpOwner) {
			pos.Equals(testEtp)
		}
	})
}

func TestKeeper_RemoveCCC(t *testing.T) {
	t.Run("Non existing ETP returns error", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		_, err := k.RemoveCCC(ctx, testEtpOwner, "notExists", *testEtp.Credits)
		errMsg := fmt.Sprintf("position for user with address %s and id %s does not exist", testEtpOwner, "notExists")
		require.Error(t, err)
		require.Equal(t, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg).Error(), err.Error())
	})

	t.Run("Existing ETP can't be modified before freeze period passes", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		params := validParams
		params.FreezePeriod = time.Minute
		assert.NotEqual(t, params.FreezePeriod, validParams.FreezePeriod)
		require.NoError(t, k.UpdateParams(ctx, params))

		k.SetPosition(ctx, testEtp)

		// require.NoError(t, k.bankKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool))

		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, *testEtp.Credits)
		require.Error(t, err)
		require.EqualError(t, err, sdkErr.Wrap(sdkErr.ErrInvalidRequest, "cannot burn position yet in the freeze period").Error())
	})

	t.Run("Existing ETP but tokens requested to burn are more than initially requested", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)

		require.NoError(t, k.bankKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool))

		require.NotZero(t, testEtp.Credits.Amount)
		burnAmountBigger := *testEtp.Credits
		burnAmountBigger = burnAmountBigger.Add(*testEtp.Credits)
		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, burnAmountBigger)
		require.Error(t, err)
	})

	t.Run("Existing ETP but cannot send tokens from from sender to module", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		_ = k.bankKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, *testEtp.Credits)
		require.Error(t, err)
	})

	t.Run("Existing ETP but cannot send collateral from module to sender", func(t *testing.T) {
		ctx, bk, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		// _ = k.bankKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		coins := sdk.NewCoins(*testEtp.Credits)
		_ = bk.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_ = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, testEtpOwner, coins)
		//_ = k.bankKeeper.AddCoins(ctx, testEtpOwner, sdk.NewCoins(*testEtp.Credits))
		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, *testEtp.Credits)
		require.Error(t, err)
		// require.Equal(t, sdk.NewInt(testEtp.Collateral), bk.GetAllBalances(ctx, testEtpOwner).AmountOf(types.BondDenom))
	})

	// TODO: control tests and remake them
	/*t.Run("Existing ETP is closed properly", func(t *testing.T) {
		ctx, bk, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)

		coins := sdk.NewCoins(*testEtp.Credits)
		_ = bk.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		_ = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, testEtpOwner, coins)
		//_ = k.bankKeeper.AddCoins(ctx, testEtpOwner, sdk.NewCoins(*testEtp.Credits))
		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, *testEtp.Credits)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(testEtp.Collateral), k.bankKeeper.GetAllBalances(ctx, testEtpOwner).AmountOf(types.BondDenom))
	})*/

	t.Run("Existing ETP returns correct residual", func(t *testing.T) {
		ctx, bk, _, k := SetupTestInput()

		k.SetPosition(ctx, testEtp)
		baseUcccAccount := sdk.NewCoin(types.CreditsDenom, sdk.NewInt(50))
		baseUcommercioAccount := sdk.NewCoin(types.BondDenom, sdk.NewInt(0))
		_ = k.bankKeeper.MintCoins(ctx, types.ModuleName, testLiquidityPool)
		coins := sdk.NewCoins(baseUcommercioAccount, baseUcccAccount)
		_ = bk.MintCoins(ctx, types.ModuleName, coins)
		_ = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, testEtpOwner, coins)
		//_ = k.bankKeeper.AddCoins(ctx, testEtpOwner, sdk.NewCoins(baseUcommercioAccount, baseUcccAccount))
		_, err := k.RemoveCCC(ctx, testEtpOwner, testEtp.ID, halfCoinSub)
		require.NoError(t, err)
		require.Equal(t, baseUcccAccount.Amount.Sub(halfCoinSub.Amount), k.bankKeeper.GetAllBalances(ctx, testEtpOwner).AmountOf(types.CreditsDenom))

		burnAmountDec := sdk.NewDecFromInt(halfCoinSub.Amount)
		collateralAmount := burnAmountDec.Mul(testEtp.ExchangeRate).Ceil().TruncateInt()

		require.Equal(t, collateralAmount, k.bankKeeper.GetAllBalances(ctx, testEtpOwner).AmountOf(types.BondDenom))
	})

}
func TestKeeper_deletePosition(t *testing.T) {

	addOneSecond := func(t time.Time) *time.Time {
		result := t.AddDate(0, 0, 1)
		return &result
	}

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
				CreatedAt:  addOneSecond(*testEtp.CreatedAt),
			},
			shouldBeDeleted: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

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

func TestKeeper_GetPositionById(t *testing.T) {

	tests := []struct {
		name     string
		id       string
		inserted []types.Position
		want     types.Position
		want1    bool
	}{
		{
			name:     "Find among only one",
			id:       testEtp.ID,
			inserted: []types.Position{testEtp},
			want:     testEtp,
			want1:    true,
		},
		{
			name:     "Find among many",
			id:       testEtp.ID,
			inserted: []types.Position{testEtp1, testEtp, testEtp2},
			want:     testEtp,
			want1:    true,
		},
		{
			name:     "Not inserted",
			id:       testEtp.ID,
			inserted: []types.Position{},
			want:     types.Position{},
			want1:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, _, _, k := SetupTestInput()

			for _, p := range tt.inserted {
				require.NoError(t, k.SetPosition(ctx, p))
			}

			got, got1 := k.GetPositionById(ctx, tt.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetPositionById() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Keeper.GetPositionById() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
