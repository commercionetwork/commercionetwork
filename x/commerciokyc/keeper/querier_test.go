package keeper_test

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/simapp"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestNewQuerier_InvalidMsg(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	app := simapp.Setup(false)
	legacyAmino := app.LegacyAmino()
	querier := keeper.NewQuerier(k, legacyAmino)
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
			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()

			for _, i := range test.storedInvites {
				k.SaveInvite(ctx, i)
			}

			querier := keeper.NewQuerier(k, legacyAmino)
			path := []string{types.QueryGetInvites}
			actualBz, _ := querier(ctx, path, request)

			var actual types.Invites
			var invites []*types.Invite
			legacyAmino.MustUnmarshalJSON(actualBz, &invites)
			for _, invite := range invites {
				actual = append(actual, *invite)
			}

			//k.Cdc.MustUnmarshalJSON(actualBz, &actual)
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
			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()

			for _, t := range test.storedTsps {
				k.AddTrustedServiceProvider(ctx, t)
			}

			querier := keeper.NewQuerier(k, legacyAmino)
			request := abci.RequestQuery{}

			path := []string{types.QueryGetTrustedServiceProviders}
			actualBz, _ := querier(ctx, path, request)

			var actual types.TrustedServiceProviders
			k.Cdc.MustUnmarshalJSON(actualBz, &actual)

			for _, tsp := range test.expected {
				require.Contains(t, actual.Addresses, tsp.String())
			}
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
			name:               "Not found membership on empty set returns correctly",
			existingMembership: types.Membership{ExpiryAt: &testExpiration},
			mustErr:            true,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()
		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		owner, _ := sdk.AccAddressFromBech32(test.existingMembership.Owner)
		tsp, _ := sdk.AccAddressFromBech32(test.existingMembership.TspAddress)
		curTime := time.Now()
		emptyMembership := types.Membership{ExpiryAt: &curTime}
		if !emptyMembership.Equals(test.existingMembership) {
			_ = k.AssignMembership(ctx, owner, test.existingMembership.MembershipType, tsp, *test.existingMembership.ExpiryAt)
		}

		querier := keeper.NewQuerier(k, legacyAmino)

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
		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()

		for _, m := range test.existingMemberships {
			owner, _ := sdk.AccAddressFromBech32(m.Owner)
			tsp, _ := sdk.AccAddressFromBech32(m.TspAddress)

			_ = k.AssignMembership(ctx, owner, m.MembershipType, tsp, *m.ExpiryAt)
		}

		querier := keeper.NewQuerier(k, legacyAmino)
		request := abci.RequestQuery{}

		path := []string{types.QueryGetMemberships}
		actualBz, _ := querier(ctx, path, request)

		var actual types.Memberships
		legacyAmino.MustUnmarshalJSON(actualBz, &actual)
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
		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()

		for _, m := range test.existingMemberships {
			owner, _ := sdk.AccAddressFromBech32(m.Owner)
			tsp, _ := sdk.AccAddressFromBech32(m.TspAddress)

			_ = k.AssignMembership(ctx, owner, m.MembershipType, tsp, *m.ExpiryAt)
		}
		k.AddTrustedServiceProvider(ctx, test.tsp)
		querier := keeper.NewQuerier(k, legacyAmino)

		path := []string{types.QueryGetTspMemberships, test.tsp.String()}
		actualBz, _ := querier(ctx, path, request)

		var actual types.Memberships
		legacyAmino.MustUnmarshalJSON(actualBz, &actual)
		require.Equal(t, test.expected, actual)

	}
}
