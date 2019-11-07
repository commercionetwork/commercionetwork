package keeper_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, test.credentials, result)
		})
	}
}

func TestKeeper_GetCredentials(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	user1, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	user2, _ := sdk.AccAddressFromBech32("cosmos1vt9vnyhukw65vvqxzp0vdnvddqlc9k88x2ajrm")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	assert.Equal(t, types.Credentials{}, k.GetCredentials(ctx))

	credential1 := types.NewCredential(user1, verifier, 0)
	k.SaveCredential(ctx, credential1)
	assert.Equal(t, types.Credentials{credential1}, k.GetCredentials(ctx))

	credential2 := types.NewCredential(user2, verifier, 0)
	k.SaveCredential(ctx, credential2)
	assert.Equal(t, types.Credentials{credential1, credential2}, k.GetCredentials(ctx))
}

func TestKeeper_SaveCredential(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	assert.Equal(t, types.Credentials{}, k.GetUserCredentials(ctx, user))

	credential1 := types.NewCredential(user, verifier, 0)
	k.SaveCredential(ctx, credential1)
	assert.Equal(t, types.Credentials{credential1}, k.GetUserCredentials(ctx, user))

	credential2 := types.NewCredential(user, verifier, 1)
	k.SaveCredential(ctx, credential2)
	assert.Equal(t, types.Credentials{credential1, credential2}, k.GetUserCredentials(ctx, user))
}
