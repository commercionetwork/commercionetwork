package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) getCredentialsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.CredentialsStorePrefix + user.String())
}

// GetUserCredentials returns the credentials that have been associated to the given user
func (keeper Keeper) GetUserCredentials(ctx sdk.Context, user sdk.AccAddress) types.Credentials {
	store := ctx.KVStore(keeper.StoreKey)

	var credentials types.Credentials
	bz := store.Get(keeper.getCredentialsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(bz, &credentials)

	return credentials
}

// GetCredentials returns the list of all the credentials that have been saved
func (keeper Keeper) GetCredentials(ctx sdk.Context) (credentials types.Credentials) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialsStorePrefix))

	credentials = types.Credentials{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &credential)
		credentials, _ = credentials.AppendIfMissing(credential)
	}

	return credentials
}

// SaveCredential allows to save the given credential inside the store
func (keeper Keeper) SaveCredential(ctx sdk.Context, credential types.Credential) {
	credentials := keeper.GetUserCredentials(ctx, credential.User)
	if credentials, edited := credentials.AppendIfMissing(credential); edited {
		store := ctx.KVStore(keeper.StoreKey)
		store.Set(keeper.getCredentialsStoreKey(credential.User), keeper.cdc.MustMarshalBinaryBare(&credentials))
	}
}
