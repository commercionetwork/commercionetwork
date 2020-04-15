package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/docs/types"
)

// UserReceivedDocumentsIterator returns an Iterator for all the received Documents of a user.
func (keeper Keeper) UserReceivedDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getReceivedDocumentsIdsStoreKey(user))
}

// UserSentDocumentsIterator returns an Iterator for all the sent Documents of a user.
func (keeper Keeper) UserSentDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getSentDocumentsIdsStoreKey(user))
}

// DocumentsIterator returns an Iterator for all the Documents saved in the store.
func (keeper Keeper) DocumentsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.DocumentStorePrefix))
}

// UserReceivedReceiptsIterator returns an Iterator for all the Document Received Receipts for a user.
func (keeper Keeper) UserReceivedReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getReceivedReceiptsIdsStoreKey(user))
}

// UserSentReceiptsIterator returns an Iterator for all the Document Sent Receipts for a user.
func (keeper Keeper) UserSentReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getSentReceiptsIdsStoreKey(user))
}

// ReceiptsIterator returns an Iterator for Sent and Received receipts.
func (keeper Keeper) ReceiptsIterators(ctx sdk.Context) (sdk.Iterator, sdk.Iterator) {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix)),
		sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
}

// SupportedMetadataSchemesIterator returns an Iterators for all the Supported Metadata Schemes.
func (keeper Keeper) SupportedMetadataSchemesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SupportedMetadataSchemesStoreKey))
}

// TrustedSchemaProposersIterator returns an Iterator for all the Trusted Schema Proposers.
func (keeper Keeper) TrustedSchemaProposersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.MetadataSchemaProposersStoreKey))
}
