package keeper_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetUserCredentials(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	tests := []struct {
		name        string
		credentials types.Credentials
	}{
		{
			name:        "Empty list returns correctly",
			credentials: types.Credentials{},
		},
		{
			name:        "Non empty list returns correctly",
			credentials: types.Credentials{types.NewCredential(user, verifier, 0)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			if len(test.credentials) > 0 {
				store.Set([]byte(types.CredentialsStorePrefix+user.String()), k.Cdc.MustMarshalBinaryBare(&test.credentials))
			}

			result := k.GetUserCredentials(ctx, user)
			require.Equal(t, test.credentials, result)
		})
	}
}

func TestKeeper_GetCredentials(t *testing.T) {
	user1, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	tests := []struct {
		name        string
		credentials types.Credentials
		expected    types.Credentials
	}{
		{
			name:        "Empty credentials returns correctly",
			credentials: types.Credentials{},
			expected:    types.Credentials{},
		},
		{
			name: "Non empty credentials returns correctly",
			credentials: types.Credentials{
				types.NewCredential(user1, verifier, 10),
				types.NewCredential(user1, verifier, 50),
			},
			expected: types.Credentials{
				types.NewCredential(user1, verifier, 10),
				types.NewCredential(user1, verifier, 50),
			},
		},
		{
			name: "Credentials with same data are not reurned",
			credentials: types.Credentials{
				types.NewCredential(user1, verifier, 10),
				types.NewCredential(user1, verifier, 10),
			},
			expected: types.Credentials{types.NewCredential(user1, verifier, 10)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, c := range test.credentials {
				k.SaveCredential(ctx, c)
			}

			require.Equal(t, test.expected, k.GetCredentials(ctx))
		})
	}
}

func TestKeeper_SaveCredential(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	tests := []struct {
		name        string
		credentials types.Credentials
		expected    types.Credentials
	}{
		{
			name:        "Empty list is returned properly",
			credentials: types.Credentials{},
			expected:    types.Credentials{},
		},
		{
			name: "Double credentials are not added",
			credentials: types.Credentials{
				types.NewCredential(user, verifier, 0),
				types.NewCredential(user, verifier, 0),
			},
			expected: types.Credentials{types.NewCredential(user, verifier, 0)},
		},
		{
			name: "Multiple credentials are returned properly",
			credentials: types.Credentials{
				types.NewCredential(user, verifier, 0),
				types.NewCredential(user, verifier, 10),
			},
			expected: types.Credentials{
				types.NewCredential(user, verifier, 0),
				types.NewCredential(user, verifier, 10),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, c := range test.credentials {
				k.SaveCredential(ctx, c)
			}

			actual := k.GetUserCredentials(ctx, user)
			require.Equal(t, test.expected, actual)
		})
	}
}
