package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

var request abci.RequestQuery

func TestNewQuerier_InvalidMsg(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)
	_, res := querier(ctx, []string{""}, abci.RequestQuery{})
	require.Error(t, res)
}

func Test_queryGetInvites(t *testing.T) {
	tests := []struct {
		name          string
		storedInvites types.Invites
		expected      types.Invites
	}{
		// These tests are not valid because can't get specific invite
		/*{
			name:          "Specific user and empty invites returns properly",
			storedInvites: types.Invites{},
			path:          []string{types.QueryGetInvites, testUser.String()},
			expected:      types.Invites{},
		},*/
		/*{
			name: "Specific user and existing invite is returned properly",
			storedInvites: types.Invites{
				types.NewInvite(testInviteSender, testUser, "bronze"),
				types.NewInvite(testInviteSender, testUser2, "bronze"),
			},
			path:     []string{types.QueryGetInvites, testUser.String()},
			expected: types.Invites{types.NewInvite(testInviteSender, testUser, "bronze")},
		},*/
		{
			name:          "All invites and empty list is returned properly",
			storedInvites: types.Invites{},
			expected:      types.Invites{},
		},
		{
			name: "All invites and non empty list is returned properly",
			storedInvites: types.Invites{
				types.NewInvite(testInviteSender, testUser, "bronze"),
				types.NewInvite(testInviteSender, testUser2, "bronze"),
			},
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

			querier := NewQuerier(k)
			path := []string{types.QueryGetInvites}
			actualBz, _ := querier(ctx, path, request)

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

			querier := NewQuerier(k)
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

			querier := NewQuerier(k)
			request := abci.RequestQuery{}

			path := []string{types.QueryGetPoolFunds}
			actualBz, _ := querier(ctx, path, request)

			var actual sdk.Coins
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)
			require.True(t, test.pool.IsEqual(actual))
		})
	}
}

func Test_queryGetMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership types.Membership
		expected           types.Membership
		mustErr            bool
	}{
		{
			name:               "Existing membership is returned properly",
			existingMembership: types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
			expected:           types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
			mustErr:            false,
		},
		{
			name:               "Not found membership returns correctly",
			existingMembership: types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			mustErr:            true,
		},
		{
			name:    "Not found membership on empty set returns correctly",
			mustErr: true,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		if !(types.Membership{}).Equals(test.existingMembership) {
			_ = k.AssignMembership(ctx, test.existingMembership.Owner, test.existingMembership.MembershipType, test.existingMembership.TspAddress, test.existingMembership.ExpiryAt)
		}

		querier := NewQuerier(k)

		path := []string{types.QueryGetMembership, testUser.String()}
		actualBz, err := querier(ctx, path, request)

		if !test.mustErr {
			require.NoError(t, err)
			var actual types.Membership
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)
			require.Equal(t, test.expected, actual)
		} else {
			require.Error(t, err)
		}
	}
}

func Test_queryGetMemberships(t *testing.T) {
	tests := []struct {
		name                string
		existingMemberships types.Memberships
		expected            types.Memberships
		mustErr             bool
	}{
		{
			name: "Existing memberships is returned properly",
			existingMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, testExpiration),
			},
			expected: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, testExpiration),
			},
		},
		{
			name:                "Not found membership returns correctly",
			existingMemberships: types.Memberships{},
			expected:            types.Memberships(nil), //TODO FIX THIS: should be types.Memberships{}
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		for _, m := range test.existingMemberships {
			_ = k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
		}

		querier := NewQuerier(k)
		request := abci.RequestQuery{}

		path := []string{types.QueryGetMemberships}
		actualBz, _ := querier(ctx, path, request)

		var actual types.Memberships
		k.Cdc.MustUnmarshalJSON(actualBz, &actual)
		require.Equal(t, test.expected, actual)

	}
}

func Test_queryGetTspMemberships(t *testing.T) {
	tests := []struct {
		name                string
		existingMemberships types.Memberships
		tsp                 sdk.AccAddress
		expected            types.Memberships
		mustErr             bool
	}{
		{
			name: "All memberships for tsp is returned properly",
			existingMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, testExpiration),
			},
			tsp: testTsp,
			expected: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, testExpiration),
			},
		},
		{
			name: "Existing memberships for tsp is returned properly",
			existingMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testUser, testExpiration),
			},
			tsp: testTsp,
			expected: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testTsp, testExpiration),
			},
		},
		{
			name: "Not found memberships for tsp returns correctly",
			existingMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser, testUser2, testExpiration),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testUser, testExpiration),
			},
			tsp:      testTsp,
			expected: types.Memberships(nil), //TODO FIX THIS: should be types.Memberships{}
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		for _, m := range test.existingMemberships {
			_ = k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
		}
		k.AddTrustedServiceProvider(ctx, test.tsp)
		querier := NewQuerier(k)

		path := []string{types.QueryGetTspMemberships, test.tsp.String()}
		actualBz, _ := querier(ctx, path, request)

		var actual types.Memberships
		k.Cdc.MustUnmarshalJSON(actualBz, &actual)
		require.Equal(t, test.expected, actual)

	}
}
