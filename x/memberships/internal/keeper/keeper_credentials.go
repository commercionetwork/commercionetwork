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
	store := ctx.KVStore(k.storeKey)

	var credentials types.Credentials
	bz := store.Get(k.getCredentialsStoreKey(user))
	k.cdc.MustUnmarshalBinaryBare(bz, &credentials)

	return credentials
}

// GetCredentials returns the list of all the credentials that have been saved
func (k Keeper) GetCredentials(ctx sdk.Context) (credentials types.Credentials) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialsStorePrefix))

	credentials = types.Credentials{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &credential)
		credentials, _ = credentials.AppendIfMissing(credential)
	}

	return credentials
}

// SaveCredential allows to save the given credential inside the store
func (k Keeper) SaveCredential(ctx sdk.Context, credential types.Credential) {
	credentials := k.GetUserCredentials(ctx, credential.User)
	if credentials, edited := credentials.AppendIfMissing(credential); edited {
		store := ctx.KVStore(k.storeKey)
		store.Set(k.getCredentialsStoreKey(credential.User), k.cdc.MustMarshalBinaryBare(&credentials))
	}
}
