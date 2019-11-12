package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getCredentialsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.CredentialsStorePrefix + user.String())
}

// GetUserCredentials returns the credentials that have been associated to the given user
func (k Keeper) GetUserCredentials(ctx sdk.Context, user sdk.AccAddress) types.Credentials {
	store := ctx.KVStore(k.StoreKey)
	key := k.getCredentialsStoreKey(user)

	credentials := types.Credentials{}
	if store.Has(key) {
		k.Cdc.MustUnmarshalBinaryBare(store.Get(key), &credentials)
	}

	return credentials
}

// GetCredentials returns the list of all the credentials that have been saved
func (k Keeper) GetCredentials(ctx sdk.Context) (credentials types.Credentials) {
	store := ctx.KVStore(k.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialsStorePrefix))

	credentials = types.Credentials{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var userCredentials types.Credentials
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &userCredentials)

		for _, userCredential := range userCredentials {
			credentials, _ = credentials.AppendIfMissing(userCredential)
		}
	}

	return credentials
}

// SaveCredential allows to save the given credential inside the store
func (k Keeper) SaveCredential(ctx sdk.Context, credential types.Credential) {
	credentials := k.GetUserCredentials(ctx, credential.User)
	if credentials, edited := credentials.AppendIfMissing(credential); edited {
		store := ctx.KVStore(k.StoreKey)
		store.Set(k.getCredentialsStoreKey(credential.User), k.Cdc.MustMarshalBinaryBare(&credentials))
	}
}
