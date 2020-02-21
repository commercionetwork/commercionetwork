package keeper_test

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestNewQuerier_InvalidMsg(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := keeper.NewQuerier(k)
	_, res := querier(ctx, []string{""}, abci.RequestQuery{})
	require.Error(t, res)
}

func Test_queryGetInvites(t *testing.T) {
	tests := []struct {
		name          string
		storedInvites types.Invites
		path          []string
		expected      types.Invites
	}{
		{
			name:          "Specific user and empty invites returns properly",
			storedInvites: types.Invites{},
			path:          []string{types.QueryGetInvites, testUser.String()},
			expected:      types.Invites{},
		},
		{
			name: "Specific user and existing invite is returned properly",
			storedInvites: types.Invites{
				types.NewInvite(testInviteSender, testUser, "bronze"),
				types.NewInvite(testInviteSender, testUser2, "bronze"),
			},
			path:     []string{types.QueryGetInvites, testUser.String()},
			expected: types.Invites{types.NewInvite(testInviteSender, testUser, "bronze")},
		},
		{
			name:          "All invites and empty list is returned properly",
			storedInvites: types.Invites{},
			path:          []string{types.QueryGetInvites},
			expected:      types.Invites{},
		},
		{
			name: "All invites and non empty list is returned properly",
			storedInvites: types.Invites{
				types.NewInvite(testInviteSender, testUser, "bronze"),
				types.NewInvite(testInviteSender, testUser2, "bronze"),
			},
			path: []string{types.QueryGetInvites},
			expected: types.Invites{
				types.NewInvite(testInviteSender, testUser2, "bronze"),
				types.NewInvite(testInviteSender, testUser, "bronze"),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, i := range test.storedInvites {
				k.SaveInvite(ctx, i)
			}

			querier := keeper.NewQuerier(k)
			actualBz, _ := querier(ctx, test.path, request)

			var actual types.Invites
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)
			require.True(t, test.expected.Equals(actual))
		})
	}

}

func Test_queryGetSigners(t *testing.T) {
	tests := []struct {
		name       string
		storedTsps ctypes.Addresses
		expected   ctypes.Addresses
	}{
		{
			name:       "Empty list is returned properly",
			storedTsps: ctypes.Addresses{},
			expected:   ctypes.Addresses{},
		},
		{
			name:       "Existing list is returned properly",
			storedTsps: []sdk.AccAddress{testUser, testTsp},
			expected:   ctypes.Addresses{testUser, testTsp},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, t := range test.storedTsps {
				k.AddTrustedServiceProvider(ctx, t)
			}

			querier := keeper.NewQuerier(k)
			request := abci.RequestQuery{}

			path := []string{types.QueryGetTrustedServiceProviders}
			actualBz, _ := querier(ctx, path, request)

			var actual []sdk.AccAddress
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)

			for _, tsp := range test.expected {
				require.Contains(t, actual, tsp)
			}
		})
	}
}

func Test_queryGetPoolFunds(t *testing.T) {
	tests := []struct {
		name string
		pool sdk.Coins
	}{
		{
			name: "Empty pool is returned properly",
			pool: sdk.Coins{},
		},
		{
			name: "Exiting pool is returned properly",
			pool: sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
				sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
			),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			if !test.pool.Empty() {
				_ = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, test.pool)
			}

			querier := keeper.NewQuerier(k)
			request := abci.RequestQuery{}

			path := []string{types.QueryGetPoolFunds}
			actualBz, _ := querier(ctx, path, request)

			var actual sdk.Coins
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)
			require.True(t, test.pool.IsEqual(actual))
		})
	}
}

func Test_queryResolveMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership types.Membership
		expected           keeper.MembershipResult
		mustErr            bool
	}{
		{
			name:               "Existing membership is returned properly",
			existingMembership: types.NewMembership(types.MembershipTypeGold, testUser),
			expected:           keeper.MembershipResult{User: testUser, MembershipType: types.MembershipTypeGold},
			mustErr:            false,
		},
		{
			name:     "Not found membership returns correctly",
			expected: keeper.MembershipResult{User: testUser, MembershipType: ""},
			mustErr:  true,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		if !(types.Membership{}).Equals(test.existingMembership) {
			_ = k.AssignMembership(ctx, test.existingMembership.Owner, test.existingMembership.MembershipType)
		}

		querier := keeper.NewQuerier(k)

		path := []string{types.QueryGetMembership, testUser.String()}
		actualBz, err := querier(ctx, path, request)

		if !test.mustErr {
			require.NoError(t, err)
			var actual keeper.MembershipResult
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)
			require.Equal(t, test.expected, actual)
		} else {
			require.Error(t, err)
		}
	}
}
