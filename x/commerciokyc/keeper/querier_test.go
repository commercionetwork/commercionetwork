package keeper

import (
	"testing"

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
	querier := NewQuerier(k, legacyAmino)
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

			querier := NewQuerier(k, legacyAmino)
			path := []string{types.QueryGetInvites}
			actualBz, _ := querier(ctx, path, request)

			var actual types.Invites
			var invites []*types.Invite
			legacyAmino.MustUnmarshalJSON(actualBz, &invites)
			for _, invite := range invites {
				actual = append(actual, *invite)
			}

			//k.cdc.MustUnmarshalJSON(actualBz, &actual)
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

			querier := NewQuerier(k, legacyAmino)
			request := abci.RequestQuery{}

			path := []string{types.QueryGetTrustedServiceProviders}
			actualBz, _ := querier(ctx, path, request)

			var actual types.TrustedServiceProviders
			k.cdc.MustUnmarshalJSON(actualBz, &actual)

			for _, tsp := range test.expected {
				require.Contains(t, actual.Addresses, tsp.String())
			}
		})
	}
}
